package utilities

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// RemoveGoogleAPI removes Google API imports and options from proto files
func RemoveGoogleAPI(ctx context.Context, compiler *protocompile.Compiler, files []string, outputDir string) (map[string][]byte, error) {
	fileDetails := make(map[string][]byte)
	fds, err := compiler.Compile(ctx, files...)
	if err != nil {
		return nil, err
	}

	for _, fd := range fds {
		var builder strings.Builder

		if fd.Syntax() != protoreflect.Proto3 {
			continue
		}
		// Write syntax
		builder.WriteString(fmt.Sprintf("syntax = \"%s\";\n\n", fd.Syntax()))

		// Write package
		if fd.Package() != "" {
			builder.WriteString(fmt.Sprintf("package %s;\n\n", fd.Package()))
		}

		// Write imports excluding Google API imports and google/protobuf/descriptor.proto
		imports := fd.Imports()
		hasImports := false
		for i := 0; i < imports.Len(); i++ {
			imp := imports.Get(i)
			importPath := string(imp.Path())
			// Skip google/api imports and google/protobuf/descriptor.proto
			if !strings.HasPrefix(importPath, "google/api") &&
				!strings.Contains(importPath, "google/protobuf/descriptor.proto") {
				builder.WriteString(fmt.Sprintf("import \"%s\";\n", importPath))
				hasImports = true
			}
		}
		if hasImports {
			builder.WriteString("\n")
		}

		// Write file options exactly as they are (no modifications)
		opts := fd.Options().(*descriptorpb.FileOptions)
		if opts != nil {
			if opts.GoPackage != nil {
				builder.WriteString(fmt.Sprintf("option go_package = \"%s\";\n", *opts.GoPackage))
			}
			if opts.JavaMultipleFiles != nil {
				builder.WriteString(fmt.Sprintf("option java_multiple_files = %t;\n", *opts.JavaMultipleFiles))
			}
			if opts.JavaOuterClassname != nil {
				builder.WriteString(fmt.Sprintf("option java_outer_classname = \"%s\";\n", *opts.JavaOuterClassname))
			}
			if opts.JavaPackage != nil {
				builder.WriteString(fmt.Sprintf("option java_package = \"%s\";\n", *opts.JavaPackage))
			}
		}
		builder.WriteString("\n")

		// Write enums
		enums := fd.Enums()
		for i := 0; i < enums.Len(); i++ {
			enum := enums.Get(i)
			writeEnum(&builder, enum)
		}

		// Write messages
		messages := fd.Messages()
		for i := 0; i < messages.Len(); i++ {
			msg := messages.Get(i)
			writeMessage(&builder, msg, 0)
		}

		// Write services - special handling for google/longrunning
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			svc := services.Get(i)
			writeService(&builder, svc, fd)
		}

		// Store the built proto content in map with key as outputDir + fd.Path()
		key := filepath.Join(outputDir, fd.Path())
		fileDetails[key] = []byte(builder.String())
	}

	return fileDetails, nil
}

// writeEnum writes an enum definition to the builder
func writeEnum(builder *strings.Builder, enum protoreflect.EnumDescriptor) {
	builder.WriteString(fmt.Sprintf("enum %s {\n", enum.Name()))
	values := enum.Values()
	for i := 0; i < values.Len(); i++ {
		value := values.Get(i)
		builder.WriteString(fmt.Sprintf("  %s = %d;\n", value.Name(), value.Number()))
	}
	builder.WriteString("}\n\n")
}

// writeMessage writes a message definition to the builder with proper indentation
func writeMessage(builder *strings.Builder, msg protoreflect.MessageDescriptor, indent int) {
	indentStr := strings.Repeat("  ", indent)

	builder.WriteString(fmt.Sprintf("%smessage %s {\n", indentStr, msg.Name()))

	// Write nested enums
	nestedEnums := msg.Enums()
	for i := 0; i < nestedEnums.Len(); i++ {
		nestedEnum := nestedEnums.Get(i)
		builder.WriteString(fmt.Sprintf("%s  enum %s {\n", indentStr, nestedEnum.Name()))
		values := nestedEnum.Values()
		for j := 0; j < values.Len(); j++ {
			value := values.Get(j)
			builder.WriteString(fmt.Sprintf("%s    %s = %d;\n", indentStr, value.Name(), value.Number()))
		}
		builder.WriteString(fmt.Sprintf("%s  }\n\n", indentStr))
	}

	// Write oneofs
	oneofs := msg.Oneofs()
	for i := 0; i < oneofs.Len(); i++ {
		oneof := oneofs.Get(i)
		builder.WriteString(fmt.Sprintf("%s  oneof %s {\n", indentStr, oneof.Name()))
		fields := msg.Fields()
		for j := 0; j < fields.Len(); j++ {
			field := fields.Get(j)
			if field.ContainingOneof() == oneof {
				writeField(builder, field, indent+2)
			}
		}
		builder.WriteString(fmt.Sprintf("%s  }\n", indentStr))
	}

	// Write regular fields (not in oneofs)
	fields := msg.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if field.ContainingOneof() == nil {
			writeField(builder, field, indent+1)
		}
	}

	// Write nested messages
	nestedMessages := msg.Messages()
	for i := 0; i < nestedMessages.Len(); i++ {
		nested := nestedMessages.Get(i)
		writeMessage(builder, nested, indent+1)
	}

	builder.WriteString(fmt.Sprintf("%s}\n\n", indentStr))
}

// writeField writes a field definition to the builder
func writeField(builder *strings.Builder, field protoreflect.FieldDescriptor, indent int) {
	indentStr := strings.Repeat("  ", indent)
	fieldType := getFieldType(field)
	builder.WriteString(fmt.Sprintf("%s%s %s = %d;\n",
		indentStr, fieldType, field.Name(), field.Number()))
}

// writeService writes a service definition to the builder
func writeService(builder *strings.Builder, svc protoreflect.ServiceDescriptor, fd protoreflect.FileDescriptor) {
	builder.WriteString(fmt.Sprintf("service %s {\n", svc.Name()))
	methods := svc.Methods()
	for i := 0; i < methods.Len(); i++ {
		method := methods.Get(i)
		writeMethod(builder, method)
	}
	builder.WriteString("}\n\n")
}

// writeMethod writes a method definition to the builder
func writeMethod(builder *strings.Builder, method protoreflect.MethodDescriptor) {
	inputType := getFullMessageName(method.Input())
	outputType := getFullMessageName(method.Output())

	rpcLine := fmt.Sprintf("  rpc %s(%s) returns (%s)", method.Name(), inputType, outputType)

	// Check if method has options
	opts := method.Options().(*descriptorpb.MethodOptions)
	if opts != nil && len(opts.ProtoReflect().GetUnknown()) > 0 {
		rpcLine += " {\n"
		// Here you could add logic to write non-google.api options if needed
		rpcLine += "  }"
	} else {
		rpcLine += ";"
	}

	builder.WriteString(rpcLine + "\n")
}

// getFieldType returns the string representation of a field's type,
// including repeated, message, and enum types with fully qualified names if necessary.
func getFieldType(field protoreflect.FieldDescriptor) string {
	repeatPrefix := ""
	if field.Cardinality() == protoreflect.Repeated {
		repeatPrefix = "repeated "
	}

	var baseType string
	switch field.Kind() {
	case protoreflect.StringKind:
		baseType = "string"
	case protoreflect.Int32Kind:
		baseType = "int32"
	case protoreflect.Int64Kind:
		baseType = "int64"
	case protoreflect.Uint32Kind:
		baseType = "uint32"
	case protoreflect.Uint64Kind:
		baseType = "uint64"
	case protoreflect.BoolKind:
		baseType = "bool"
	case protoreflect.FloatKind:
		baseType = "float"
	case protoreflect.DoubleKind:
		baseType = "double"
	case protoreflect.BytesKind:
		baseType = "bytes"
	case protoreflect.MessageKind:
		msgDesc := field.Message()
		baseType = getFullMessageName(msgDesc)
	case protoreflect.EnumKind:
		enumDesc := field.Enum()
		baseType = getFullEnumName(enumDesc)
	default:
		baseType = "string"
	}

	return repeatPrefix + baseType
}

// getFullMessageName returns the fully qualified name of a message, including nested paths and package.
func getFullMessageName(msgDesc protoreflect.MessageDescriptor) string {
	var parts []string
	current := msgDesc
	for current != nil {
		parts = append([]string{string(current.Name())}, parts...)
		parent := current.Parent()
		if parent == nil {
			break
		}
		if parentMsg, ok := parent.(protoreflect.MessageDescriptor); ok {
			current = parentMsg
		} else {
			break
		}
	}

	fullName := strings.Join(parts, ".")
	if pkg := msgDesc.ParentFile().Package(); pkg != "" {
		fullName = string(pkg) + "." + fullName
	}
	return fullName
}

// getFullEnumName returns the fully qualified name of an enum, including nested paths and package.
func getFullEnumName(enumDesc protoreflect.EnumDescriptor) string {
	var parts []string
	parts = append(parts, string(enumDesc.Name()))

	if parent := enumDesc.Parent(); parent != nil {
		if parentMsg, ok := parent.(protoreflect.MessageDescriptor); ok {
			parentName := getFullMessageName(parentMsg)
			return parentName + "." + strings.Join(parts, ".")
		}
	}

	fullName := strings.Join(parts, ".")
	if pkg := enumDesc.ParentFile().Package(); pkg != "" {
		fullName = string(pkg) + "." + fullName
	}
	return fullName
}
