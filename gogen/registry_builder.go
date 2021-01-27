package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

type RegistryBuilderRequest struct {
	RegistryName   string
	UsecaseName    string
	GatewayName    string
	ControllerName string
	FolderPath     string
}

type registryBuilder struct {
	RegistryBuilderRequest RegistryBuilderRequest
}

func NewRegistry(req RegistryBuilderRequest) Generator {
	return &registryBuilder{RegistryBuilderRequest: req}
}

func (d *registryBuilder) Generate() error {

	registryName := strings.TrimSpace(d.RegistryBuilderRequest.RegistryName)
	usecaseName := strings.TrimSpace(d.RegistryBuilderRequest.UsecaseName)
	gatewayName := strings.TrimSpace(d.RegistryBuilderRequest.GatewayName)
	controllerName := strings.TrimSpace(d.RegistryBuilderRequest.ControllerName)
	folderPath := d.RegistryBuilderRequest.FolderPath

	packagePath := GetPackagePath()

	if len(registryName) == 0 {
		return fmt.Errorf("Registry name must not empty")
	}

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	if len(gatewayName) == 0 {
		return fmt.Errorf("Gateway name must not empty")
	}

	if len(controllerName) == 0 {
		return fmt.Errorf("Controller name must not empty")
	}

	if !IsExist(fmt.Sprintf("%s/controller/%s/%s.go", folderPath, strings.ToLower(controllerName), PascalCase(usecaseName))) {
		return fmt.Errorf("controller %s/%s is not found", controllerName, PascalCase(usecaseName))
	}

	if !IsExist(fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, strings.ToLower(gatewayName), PascalCase(usecaseName))) {
		return fmt.Errorf("gateway %s/%s is not found", gatewayName, PascalCase(usecaseName))
	}

	if !IsExist(fmt.Sprintf("%s/usecase/%s", folderPath, strings.ToLower(usecaseName))) {
		return fmt.Errorf("usecase %s is not found", PascalCase(usecaseName))
	}

	CreateFolder("%s/application/registry", folderPath)

	CreateFolder("%s/application/router", folderPath)

	CreateFolder("%s/application/infrastructure", folderPath)

	rg := StructureRegistry{
		RegistryName: registryName,
		PackagePath:  packagePath,
	}

	_ = WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%s/main.go", folderPath),
		rg,
	)

	_ = WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%s/application/application.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/infrastructure/gracefully_shutdown._go",
		fmt.Sprintf("%s/application/infrastructure/gracefully_shutdown.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/infrastructure/http_handler._go",
		fmt.Sprintf("%s/application/infrastructure/http_handler.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/infrastructure/infrastructure._go",
		fmt.Sprintf("%s/application/infrastructure/infrastructure.go", folderPath),
		struct{}{},
	)

	// registry
	{

		_ = WriteFileIfNotExist(
			"application/registry/registry._go",
			fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)),
			rg,
		)

		funcCallInjectedCode, _ := PrintTemplate("application/registry/usecase_assign._go", d.RegistryBuilderRequest)

		registryFile := fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName))

		fSet := token.NewFileSet()
		node, errParse := parser.ParseFile(fSet, registryFile, nil, parser.ParseComments)
		if errParse != nil {
			return errParse
		}

		existingImportMap := ReadImports(node)

		file, _ := os.Open(registryFile)
		defer file.Close()
		scanner := bufio.NewScanner(file)

		methodCallMode := false
		importMode := false
		var buffer bytes.Buffer
		for scanner.Scan() {
			row := scanner.Text()

			if methodCallMode {

				if strings.HasPrefix(strings.TrimSpace(row), "return") {
					methodCallMode = false
					buffer.WriteString(funcCallInjectedCode)
					buffer.WriteString("\n")
				}

			}

			if strings.HasPrefix(row, fmt.Sprintf("func (r *%sRegistry) RegisterUsecase() map[string]interface{} {", CamelCase(registryName))) {
				methodCallMode = true

			} else //

			if importMode && strings.HasPrefix(row, ")") {
				importMode = false

				// if _, exist := existingImportMap[fmt.Sprintf("\"%s/controller/%s\"", packagePath, controllerName)]; !exist {
				// 	buffer.WriteString(fmt.Sprintf("	\"%s/controller/%s\"", packagePath, controllerName))
				// 	buffer.WriteString("\n")
				// }

				if _, exist := existingImportMap[fmt.Sprintf("\"%s/gateway/%s\"", packagePath, gatewayName)]; !exist {
					buffer.WriteString(fmt.Sprintf("	\"%s/gateway/%s\"", packagePath, gatewayName))
					buffer.WriteString("\n")
				}

				if _, exist := existingImportMap[fmt.Sprintf("\"%s/usecase/%s\"", packagePath, LowerCase(usecaseName))]; !exist {
					buffer.WriteString(fmt.Sprintf("	\"%s/usecase/%s\"", packagePath, LowerCase(usecaseName)))
					buffer.WriteString("\n")
				}

			} else //

			if strings.HasPrefix(row, "import (") {
				importMode = true

			}

			buffer.WriteString(row)
			buffer.WriteString("\n")
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)), buffer.Bytes(), 0644); err != nil {
			return err
		}

	}

	// router
	{

		rg := StructureRouter{
			ControllerName: controllerName,
			PackagePath:    packagePath,
			UsecaseName:    usecaseName,
		}

		_ = WriteFileIfNotExist(
			"application/router/router._go",
			fmt.Sprintf("%s/application/router/router.go", folderPath),
			rg,
		)

		funcCallInjectedCode, _ := PrintTemplate("application/router/controller_usecase._go", d.RegistryBuilderRequest)

		routerFile := fmt.Sprintf("%s/application/router/router.go", folderPath)

		file, _ := os.Open(routerFile)
		defer file.Close()
		scanner := bufio.NewScanner(file)

		methodCallMode := false
		var buffer bytes.Buffer
		for scanner.Scan() {
			row := scanner.Text()

			if methodCallMode {

				if strings.HasPrefix(strings.TrimSpace(row), "}") {
					methodCallMode = false
					buffer.WriteString(funcCallInjectedCode)
					buffer.WriteString("\n")
				}

			}

			if strings.HasPrefix(row, "func (m *MyRouter) RegisterRouter(usecaseMap func(string)interface{}) {") {
				methodCallMode = true

			}

			buffer.WriteString(row)
			buffer.WriteString("\n")
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/application/router/router.go", folderPath), buffer.Bytes(), 0644); err != nil {
			return err
		}

	}

	return nil
}
