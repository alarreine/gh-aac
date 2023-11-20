gh-aac
├── config
│   ├── export
│   │   └── --org <org-name> --path <file-path> --format <format>
│   └── import
│       └── --org <org-name> --path <file-path>
└── permissions
    ├── org
    │   ├── member
    │   │   ├── modify
    │   │   │   └── --org <org-name> --user <username> --role <role>
    │   │   └── remove
    │   │       └── --org <org-name> --user <username>
    │   └── team
    │       ├── modify
    │       │   └── --org <org-name> --team <team-name> --role <role>
    │       └── remove
    │           └── --org <org-name> --team <team-name>
    ├── repo
    │   ├── member
    │   │   ├── modify
    │   │   │   └── --org <org-name> --repo <repo-name> --user <username> --permissions <permissions>
    │   │   └── remove
    │   │       └── --org <org-name> --repo <repo-name> --user <username>
    │   └── team
    │       ├── modify
    │       │   └── --org <org-name> --repo <repo-name> --team <team-name> --permissions <permissions>
    │       └── remove
    │           └── --org <org-name> --repo <repo-name> --team <team-name>
    └── team
        ├── member
        │   ├── modify
        │   │   └── --org <org-name> --team <team-name> --user <username> --role <role>
        │   └── remove
        │       └── --org <org-name> --team <team-name> --user <username>
        └── child
            ├── add
            │   └── --org <org-name> --team <team-name> --child <child-team-name>
            └── remove
                └── --org <org-name> --team <team-name> --child <child-team-name>
