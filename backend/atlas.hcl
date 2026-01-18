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

env "dev" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  url = local.envfile["DEV_MIGRATION_DB_CONNECTION_STRING"]
  schemas = ["public"]
  migration {
    dir = "file://internal/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}


env "prod" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  url = local.envfile["PROD_MIGRATION_DB_CONNECTION_STRING"]
  schemas = ["public"]
  migration {
    dir = "file://internal/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}