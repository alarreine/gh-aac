# GitHub Access As Code (gh-aac)

`gh-aac` es una herramienta CLI (Interfaz de Línea de Comandos) creada para gestionar y versionar los accesos y permisos en organizaciones y repositorios de GitHub de manera estructurada y eficiente.

## Características

- Inicialización de configuraciones de organización en GitHub.
- Exportación e importación de configuraciones de acceso "as code".
- Gestión de miembros y equipos en repositorios y organizaciones.
- Versionado de todos los cambios en la configuración de acceso.
- Auditoría y revisión de los cambios en los accesos y permisos.

## Requisitos

- Go 1.16 o superior.
- [GitHub Personal Access Token](https://github.com/settings/tokens) con los permisos necesarios.
- Git

## Instalación

```bash
go get github.com/<tu-usuario>/gh-aac
```

# Uso
## Inicialización
```bash
gh-aac init --org <org-name>
```
## Exportar Configuración de Acceso
```bash
gh-aac export --org <org-name> --file <file-path>
```
## Importar Cambios
```bash
gh-aac import --org <org-name> --file <file-path>
```
## Agregar Nuevo Usuario a Repositorio
```bash
gh-aac repo add-member --org <org-name> --repo <repo-name> --user <username> --role <role>
```
## Remover Usuario de Repositorio
```bash
gh-aac repo remove-member --org <org-name> --repo <repo-name> --user <username>
```
## ... y otros comandos.

Contribución
Si deseas contribuir al proyecto, por favor lee CONTRIBUTING.md y siente libre de hacer fork del repositorio y enviar un Pull Request.

Licencia
MIT License (si tienes una licencia, sino puedes omitir esta sección o añadir la que corresponda)

Contacto
GitHub: alarreine
Email: alarreine@gmail.com