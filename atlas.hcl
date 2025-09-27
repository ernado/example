data "external_schema" "ent" {
  program = [
    "go", "tool", "ent", "schema",
    "./internal/ent/schema",
    "--dialect", "postgres",
  ]
}

env "dev" {
  dev = "docker://postgres/18/test?search_path=public"
  src  = data.external_schema.ent.url
}
