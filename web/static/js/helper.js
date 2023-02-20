const toastOptions = {
    duration: 5000,
    close: true,
    gravity: "top", // `top` or `bottom`
    position: "right", // `left`, `center` or `right`
    stopOnFocus: true, // Prevents dismissing of toast on hover
}
const helper = {
    alert: {
        success(text) {
            return Swal.fire({
                title: 'Success',
                text: text,
                icon: 'success'
            })
        },
        error(text) {
            return Swal.fire({
                title: 'Error',
                text: text,
                icon: 'error'
            })
        },
        confirm(text) {
            return Swal.fire({
                title: 'Confirmation',
                icon: 'question',
                text: text,
                confirmButtonText: 'Confirm',
                cancelButtonText: 'Cancel',
                showCancelButton: true
            })
        }
    },
    toast: {
        simple(text) {
            Toastify({
                text,
                ...toastOptions,
                style: {
                    background: "black",
                },
            }).showToast()
        },
        error(text) {
            Toastify({
                text,
                ...toastOptions,
                style: {
                    background: "red",
                },
            }).showToast()
        }
    },
    requestURI: new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    }),
}
