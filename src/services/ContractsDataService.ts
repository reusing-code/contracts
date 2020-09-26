import http from '@/http-common'

class ContractsDataService {
  getAll() {
    return http.get('/contracts')
  }
  get(id: string) {
    return http.get(`/contracts/${id}`)
  }
  create(data: any) {
    return http.post('/tutorials', data)
  }
  update(id: string, data: any) {
    return http.put(`/contracts/${id}`, data)
  }
  delete(id: string) {
    return http.delete(`/contracts/${id}`)
  }
}
