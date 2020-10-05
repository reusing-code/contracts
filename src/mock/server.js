import { Model, Server } from "miragejs";

export function makeServer({ environment = "development" } = {}) {
  let server = new Server({
    environment,

    models: {
      contract: Model,
    },

    seeds(server) {
      server.create("contract", { id: 20 });
      server.create("contract", { id: 21 });
    },

    routes() {
      this.namespace = "api";

      this.get("/contracts", (schema) => {
        return schema.contracts.all();
      });
    },
  });

  return server;
}
