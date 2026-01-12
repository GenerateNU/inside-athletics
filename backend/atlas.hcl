data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres", // | postgres | sqlite | sqlserver
  ]
}

variable "envfile" {
    type    = string
    default = "./.env"
}

locals {
    envfile = {
        for line in split("\n", file(var.envfile)): split("=", line)[0] => regex("=(.*)", line)[0]
        if !startswith(line, "#") && length(split("=", line)) > 1
    }
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://file?mode=memory"
  migration {
    dir = "file://internal/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "dev" {
  url = local.envfile["DEV_MIGRATION_DB_CONNECTION_STRING"]
  migration {
    dir = "file://internal/migrations"
  }
}

env "prod" {
  url = local.envfile["PROD_MIGRATION_DB_CONNECTION_STRING"]
  migration {
    dir = "file://internal/migrations"
    revisions_schema = "atlas_schema_revisions"
  }
}