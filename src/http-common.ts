import axios from 'axios'

export default axios.create({
  baseURL:
      (process.env.VUE_APP_API_ENDPOINT ? process.env.VUE_APP_API_ENDPOINT :
                                          '') +
      '/api/',
  headers: {'Content-type': 'application/json'}
});