const baseUrl = "https://auth.fxchampionship.com"
const urlReviews = baseUrl + "/auth/uuser/in-review"
const urlUpdateUser = baseUrl + "/auth/uuser/update-user"

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
    const userInfo = JSON.parse(localStorage.getItem("user"))
    $("#userinfo_name").html(userInfo.name)
    $("#userinfo_email").html(userInfo.email)
    if (userInfo.image === '') {
        $("#avatarImage").attr("src", "src/assets/images/profile/user-1.jpg")
    } else {
        $("#avatarImage").attr("src", userInfo.image)
    }
    if (userInfo.inreview === 'not_yet') {
        $("#indentify_status").addClass('text-danger').html(`<i class="ti ti-ban"></i> Account has not been verified!.`)
    } else {
        $("#indentify_status").removeClass('text-danger').addClass('text-success').html(`<i class="ti ti-circle-check-filled"></i>Account has been verified!`)
        $("#indentify_update").addClass("disabled")
    }
    if (userInfo.description == '') {
        $("#userinfo_description").html(`<p id="existing_description" class="card-text">Please introduce something about yourself!<br>Click here to edit.</p>`)
    } else {
        $("#userinfo_description").html(`<p id="existing_description" class="card-text">${userInfo.description}</p>`)
    }


    // $("#userinfo_description").attr("placeholder", userInfo.description)
    if (localStorage.getItem("indentify") !== null) {
        const indentify_data = JSON.parse(localStorage.getItem("indentify"))
        $("#front_selectedImage").attr("src", indentify_data.image_front)
        $("#front_selectedImage").addClass("indentify-img w-100 mt-3 rounded border border-secondary")
        $("#back_selectedImage").attr("src", indentify_data.image_back)
        $("#back_selectedImage").addClass("indentify-img w-100 mt-3 rounded border border-secondary")
    }
})

$(document).ready(function () {
    var img_front = ""
    var img_back = ""
    $("#front_formFile").on("change", function (event) {
        var file = event.target.files[0];
        var selectedImage = $("#front_selectedImage");

        if (file) {
            var reader = new FileReader();
            reader.onload = function (e) {
                selectedImage.attr("src", e.target.result);
                selectedImage.addClass("indentify-img w-100 mt-3 rounded border border-secondary")
                img_front = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    });
    $("#back_formFile").on("change", function (event) {
        var file = event.target.files[0];
        var selectedImage = $("#back_selectedImage");
        if (file) {
            var reader = new FileReader();
            reader.onload = function (e) {
                selectedImage.attr("src", e.target.result);
                selectedImage.addClass("indentify-img w-100 mt-3 rounded border border-secondary")
                img_back = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    });
    $("#indentify_update").click(function () {
        if (img_front === '' || img_back === '') {
            window.alert("Please upload the front and back images of the national ID card first.");
            return;
        }
        const inpReview = {
            "image_front": img_front,
            "image_back": img_back
        };
        const jwtToken = getCookie("token");
        const headers = new Headers({
            'Authorization': `Bearer ${jwtToken}`
        });
        $("#indentify_update").addClass("disabled")
        fetch(urlReviews, {
            method: "POST",
            headers: headers,
            body: JSON.stringify(inpReview),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json(); // Parse the response JSON if needed
            })
            .then(dataResponse => {
                localStorage.setItem('indentify', JSON.stringify(inpReview));
                window.alert("Success!")
            })
            .catch(error => {
                console.error("Error:", error);
            });
    })
});

$(document).ready(function () {
    const userInfo = JSON.parse(localStorage.getItem("user"))
    let img_avata = userInfo.image
    let inpDescription = userInfo.description
    $("#avatarInput").on("change", function (event) {
        var file = event.target.files[0];
        var selectedImage = $("#avatarImage");
        if (file) {
            var reader = new FileReader();
            reader.onload = function (e) {
                selectedImage.attr("src", e.target.result);
                img_avata = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    });
    $("#existing_description").click(function () {
        $("#userinfo_description").html(`<textarea id="inp_description" class="form-control" rows="7" aria-label="With textarea">`)
        $("#inp_description").val(userInfo.description)
    })
    $("#user_info_updated").click(function () {
        inpDescription = $("#inp_description").val()
        if (inpDescription === undefined) {
            inpDescription = userInfo.description
        }
        const jwtToken = getCookie("token");
        const headers = new Headers({
            'Authorization': `Bearer ${jwtToken}`
        });
        const inpUserUpdate = {
            "image": img_avata,
            "description": inpDescription
        };
        console.log(inpUserUpdate)
        fetch(urlUpdateUser, {
            method: "POST",
            headers: headers,
            body: JSON.stringify(inpUserUpdate),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json(); // Parse the response JSON if needed
            })
            .then(dataResponse => {
                userInfo.image = img_avata;
                userInfo.description = inpDescription;
                const user = JSON.stringify(userInfo)
                localStorage.setItem('user', JSON.stringify(userInfo));
                window.alert("Success!")
                window.location.reload()
            })
            .catch(error => {
                console.error("Error:", error);
            });
    })
})
