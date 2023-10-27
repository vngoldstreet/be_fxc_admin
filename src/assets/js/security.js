const baseUrl = "https://auth.fxchampionship.com"
const urlReviews = baseUrl + "/auth/uuser/in-review"
const urlUpdatePassword = baseUrl + "/auth/uuser/update-password"

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
    $("#user_info_change_password").click(function () {
        const inpCPassword = $('#cpassword').val();
        const inpNPassword = $('#npassword').val();
        const confirmPassword = $("#confirm_password").val();

        const specialCharacterRegex = /[!@#$%^&*()_+{}\[\]:;<>,.?~\\-]/;
        if (inpNPassword.length >= 8 && specialCharacterRegex.test(inpNPassword)) {
            // Password is valid
            $('#npassword').removeClass('is-invalid').addClass('is-valid');
            $('#fb_npassword').removeClass('invalid-feedback').addClass('valid-feedback').text('Valid password!');
        } else {
            // Password is invalid
            $('#npassword').addClass('is-invalid');
            $('#fb_npassword').removeClass('valid-feedback').addClass('invalid-feedback').text('Password must be at least 8 characters and contain special characters');
        }

        if (inpNPassword !== confirmPassword) {
            $('#confirm_password').addClass('is-invalid');
            $("#fb_confirm_password").addClass('invalid-feedback').text("Passwords do not match!");
        } else {
            // Passwords match, you can proceed with registration
            $('#confirm_password').removeClass('is-invalid').addClass('is-valid');
            $("#fb_confirm_password").removeClass('invalid-feedback').addClass('valid-feedback').text("Valid password!");
            // Here you can make an AJAX request to your registration API
        }
        const updatePassword = {
            "c_password": inpCPassword,
            "n_password": inpNPassword,
        };
        const jwtToken = getCookie("token");
        const headers = new Headers({
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwtToken}`
        });
        console.log(updatePassword);
        fetch(urlUpdatePassword, {
            method: "POST",
            headers: headers,
            body: JSON.stringify(updatePassword),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json();
            })
            .then(responseData => {
                if (responseData.status === "failure") {
                    $("#messagefailure").addClass('text-danger')
                    $("#messagefailure").text(responseData.message)
                } else {
                    $('#cpassword').val(0);
                    $('#npassword').val(0);
                    $("#confirm_password").val(0);
                    $("#messagefailure").addClass('text-success')
                    $("#messagefailure").text(responseData.message)
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
    })
})

$(document).ready(function () {
    const userInfo = JSON.parse(localStorage.getItem("user"));
    const now = new Date();
    const currentHour = now.getHours();
    let greetingText = "";

    if (currentHour >= 18) {
        greetingText = "Good evening";
    } else if (currentHour >= 12) {
        greetingText = "Good afternoon";
    } else {
        greetingText = "Good morning";
    }
    const greetingMessage = `${greetingText}: ${userInfo.name} (${userInfo.email})!`;
    $("#username").text(greetingMessage);

    $("#submit_logout").click(function () {
        // Clear the 'token' cookie and redirect to the login page
        document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/';
        window.location.href = "/login";
    });
});