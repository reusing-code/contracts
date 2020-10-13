import { Model, Server } from "miragejs";

const mockContractData = [
  {
    id: 1,
    start: "2018-05-16",
    end: "2021-05-15",
    extensionMonths: "12",
    noticePeriodMonths: "3",
    pricePerMonth: "29.99",
    notes: "This is a note!",
    company: "Telco",
    product: "Mobile Plan",
    contractNumber: "123-555-456",
    customerNumber: "666-123-321",
    account: "dude@example.com",
    paymentOption: "VISA",
    status: "active",
  },
  {
    id: 2,
    company: "Telco",
    product: "Landline500",
    start: "2020-12-01",
    end: "2022-11-30",
    extensionMonths: "6",
    noticePeriodMonths: "1",
    pricePerMonth: "45.98",
    contractNumber: "ABC-555-DEF",
    customerNumber: "AAABBBCCC",
    account: "dude@example.com",
    paymentOption: "PayPal",
    notes: "Another note",
    status: "planned",
  },
  {
    id: 3,
    company: "Insurance Corp",
    product: "Verhicle insurance XXL",
    start: "2006-04-02",
    end: "2121-03-22",
    extensionMonths: "24",
    noticePeriodMonths: "6",
    pricePerMonth: "250",
    contractNumber: "123",
    customerNumber: "123",
    account: "dude",
    paymentOption: "CC",
    notes: "Asdf",
    status: "cancelled",
  },
  {
    id: 4,
    company: "Big Bank",
    product: "Mega Account",
    start: "1980-01-01",
    end: "1980-01-02",
    extensionMonths: "0",
    noticePeriodMonths: "0",
    pricePerMonth: "1.99",
    contractNumber: "1",
    customerNumber: "2",
    account: "123",
    paymentOption: "-",
    notes: "-",
    status: "ended",
  },
];

export function makeServer({ environment = "development" } = {}) {
  const server = new Server({
    environment,

    models: {
      contract: Model,
    },

    seeds(server) {
      mockContractData.forEach((element) => server.create("contract", element));
    },

    routes() {
      this.namespace = "api";

      this.get("/contracts", (schema) => {
        return schema.contracts.all();
      });

      this.post("/contracts", (schema, request) => {
        const data = JSON.parse(request.requestBody);
        const contracts = schema.contracts.all().models;
        const maxID = Math.max(...contracts.map((c) => c.id), 0);
        data.id = maxID + 1;
        schema.contracts.create(data);
      });

      this.delete("/contracts/:id", (schema, request) => {
        const id = request.params.id;

        return schema.contracts.find(id).destroy();
      });

      this.put("/contracts/:id", (schema, request) => {
        const id = request.params.id;

        return schema.contracts
          .find(id)
          .update(JSON.parse(request.requestBody));
      });
    },
  });

  return server;
}
