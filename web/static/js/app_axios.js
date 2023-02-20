const api = axios.create({
    baseURL: '/api/',
    timeout: 100e3,
    withCredentials: true,
})
api.interceptors.response.use(function (response) {
    if (response.data.code !== 0) {
        helper.toast.error(response.data.msg)
        return Promise.reject()
    }
    response.data = response.data.data
    return response
}, function (error) {
    if(error.response.data.msg){
        helper.toast.error(error.response.data.msg)
    }else{
        helper.toast.error('Failed to connect to the server')
    }
    return Promise.reject(error)
})
