import http from "@/http-common";

class ContractsDataService {
  getAll() {
    return http.get("/contracts");
  }
  get(id: number) {
    return http.get(`/contracts/${id}`);
  }
  create(data: any) {
    return http.post("/contracts", data);
  }
  update(data: any) {
    return http.put(`/contracts/${data.id}`, data);
  }
  delete(id: number) {
    return http.delete(`/contracts/${id}`);
  }
}

export default new ContractsDataService();
