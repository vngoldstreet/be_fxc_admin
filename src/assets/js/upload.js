let urlUpload = "/auth/upload-csv";

function getCookie(cookieName) {
    var name = cookieName + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var cookieArray = decodedCookie.split(';');

    for (var i = 0; i < cookieArray.length; i++) {
        var cookie = cookieArray[i].trim();
        if (cookie.indexOf(name) === 0) {
            return cookie.substring(name.length, cookie.length);
        }
    }
    return "";
}

$(document).ready(function () {
    $('#send_upload_file').on("click", function (e) {
        e.preventDefault()
        var fileInput = $('#inputGroupFile02')[0].files[0]; // Get the selected file
        console.log(fileInput)
        if (fileInput) {
            let jwtToken = getCookie("token");
            if (!jwtToken) {
                console.error("Error: JWT token is missing.");
                return;
            }
            var formData = new FormData();
            formData.append('file', fileInput); // Create a FormData object and append the file to it
            let headers = new Headers({
                'Authorization': `Bearer ${jwtToken}`
            });
            console.log(formData)
            fetch(urlUpload, {
                method: "POST",
                headers: headers,
                body: formData,
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error("Network response was not successful");
                    }
                    return response.json(); // Parse the response JSON if needed
                })
                .then(dataResponse => {
                    let stringData = JSON.stringify(dataResponse)
                    let html = `<code class='w-100 text-success'>${stringData}</code>`
                    $("#fb_msg_create").removeClass().addClass("fw-semibold");
                    $("#fb_msg_create").html(html);
                })
                .catch(error => {
                    let stringData = JSON.stringify(error)
                    let html = `<code class='w-100 text-danger'>${stringData}</code>`
                    $("#fb_msg_create").removeClass().addClass("fw-semibold");
                    $("#fb_msg_create").html(html);
                    console.error("Error:", error);
                });
        } else {
            alert('Please select a file to upload.');
        }
    });
});