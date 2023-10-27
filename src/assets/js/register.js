// const redirectURL = "/";
// const baseUrl = "https://auth.fxchampionship.com"
// const apiRegisterUrl = baseUrl + "/public/register";
// // Get a reference to the HTML form element
// const form = document.getElementById("formreg");
// // Add a submit event listener to the form
// form.addEventListener("submit", function (event) {
//     event.preventDefault(); // Prevent the default form submission
//     const inpName = document.getElementById('name').value;
//     const inpEmail = document.getElementById('email').value;
//     const inpPassword = document.getElementById('password').value;
//     const inpPhone = document.getElementById('phone').value;
//     const inpPasscode = document.getElementById('parnercode').value;
//     // Tạo đối tượng JSON từ giá trị thu thập
//     const inpLoginData = {
//         "email": inpEmail,
//         "password": inpPassword,
//         "name": inpName,
//         "phone": inpPhone,
//         "parnercode": inpPasscode
//     };

//     fetch(apiRegisterUrl, {
//         method: "POST",
//         headers: {
//             'Content-Type': 'application/json' // Set the content type to JSON
//         },
//         body: JSON.stringify(inpLoginData),
//     })
//         .then(response => {
//             if (!response.ok) {
//                 throw new Error("Network response was not ok");
//             }
//             return response.json(); // Parse the response JSON if needed
//         })
//         .then(responseData => {
//             console.log("Response from server:", responseData);
//             if (responseData.status == "failure") {
//                 document.getElementById("messagefailure").classList.add("text-danger")
//                 document.getElementById("messagefailure").innerHTML = responseData.message;
//             } else {
//                 document.getElementById("messagefailure").classList.add("text-success")
//                 document.getElementById("messagefailure").innerHTML = responseData.message;
//                 setTimeout((window.location.href = "/login"), 500);
//             }
//         })
//         .catch(error => {
//             console.error("Error:", error);
//         });
// });

const redirectURL = "/";
const baseUrl = "https://auth.fxchampionship.com";
const apiRegisterUrl = baseUrl + "/public/register";

$(document).ready(function () {
    $("#submit_form_signup").click(function () {
        const inpName = $('#name').val();
        const inpEmail = $('#email').val();
        const inpPassword = $('#password').val();
        const confirmPassword = $("#confirm_password").val();
        const inpPhone = $('#phone').val();
        const inpPartnerCode = $('#partnercode').val();

        if (inpName === '') {
            $('#name').addClass('is-invalid');
            $('#fb_name').addClass('invalid-feedback').text('Name is required'); // Display an error message
        } else {
            // Valid email format
            $('#name').removeClass('is-invalid').addClass('is-valid');
            $('#fb_name').removeClass('invalid-feedback').addClass('invalid-feedback').text('Look good'); // Clear the error message
        }

        const emailRegex = /^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$/; // Regular expression for email validation
        if (!inpEmail.match(emailRegex)) {
            // Invalid email format
            $('#email').addClass('is-invalid');
            $('#fb_email').addClass('invalid-feedback').text('Invalid email format'); // Display an error message
        } else {
            // Valid email format
            $('#email').removeClass('is-invalid').addClass('is-valid');
            $('#fb_email').removeClass('invalid-feedback').addClass('valid-feedback').text('Looks good!');// Clear the error message
        }
        const phonePattern = /^\+84\d{9,11}$/;
        if (phonePattern.test(inpPhone)) {
            // Valid phone number format
            $('#phone').removeClass('is-invalid').addClass('is-valid');
            $('#fb_phone').removeClass('invalid-feedback').addClass('valid-feedback').text('Looks good!');
        } else {
            // Invalid phone number format
            $('#phone').addClass('is-invalid');
            $('#fb_phone').removeClass('valid-feedback').addClass('invalid-feedback').text('Invalid Vietnamese phone number format!');
        }

        const specialCharacterRegex = /[!@#$%^&*()_+{}\[\]:;<>,.?~\\-]/;
        if (inpPassword.length >= 8 && specialCharacterRegex.test(inpPassword)) {
            // Password is valid
            $('#password').removeClass('is-invalid').addClass('is-valid');
            $('#fb_password').removeClass('invalid-feedback').addClass('valid-feedback').text('Valid password!');
        } else {
            // Password is invalid
            $('#password').addClass('is-invalid');
            $('#fb_password').removeClass('valid-feedback').addClass('invalid-feedback').text('Password must be at least 8 characters and contain special characters');
        }

        if (inpPassword !== confirmPassword) {
            $('#confirm_password').addClass('is-invalid');
            $("#fb_confirm_password").addClass('invalid-feedback').text("Passwords do not match!");
        } else {
            // Passwords match, you can proceed with registration
            $('#confirm_password').removeClass('is-invalid').addClass('is-valid');
            $("#fb_confirm_password").removeClass('invalid-feedback').addClass('valid-feedback').text("Valid password!");
            // Here you can make an AJAX request to your registration API
        }

        if (inpPartnerCode === '') {
            inpPartnerCode = "";
        }

        const registrationData = {
            "name": inpName,
            "email": inpEmail,
            "password": inpPassword,
            "phone": inpPhone,
            "partnercode": inpPartnerCode
        };

        fetch(apiRegisterUrl, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(registrationData),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json();
            })
            .then(responseData => {
                // console.log("Response from server:", responseData);
                if (responseData.status === "failure") {
                    $("#messagefailure").addClass('text-danger')
                    $("#messagefailure").text(responseData.message)
                } else {
                    $("#messagefailure").addClass('text-success')
                    $("#messagefailure").text(responseData.message)
                    setTimeout(function () {
                        window.location.href = "/login";
                    }, 500);
                }
            })
            .catch(error => {
                console.error("Error:", error);
            });
    })
})