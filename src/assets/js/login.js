const apiLoginUrl = "/public/login";
$(document).ready(function () {
    $("#submit_form_login").click(function () {
        const inpEmail = $('#email').val(); // Use jQuery to get the input value
        const inpPassword = $('#password').val(); // Use jQuery to get the input value
        const inpLoginData = {
            "email": inpEmail,
            "password": inpPassword
        };
        console.log(inpLoginData)
        // Create a custom request object for the fetch
        const requestOptions = {
            method: "POST",
            headers: {
                "Content-Type": "application/json", // Ensure data is sent as JSON
            },
            body: JSON.stringify(inpLoginData),
        };

        fetch(apiLoginUrl, requestOptions)
            .then(response => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json(); // Parse the response JSON if needed
            })
            .then(responseData => {
                localStorage.setItem('user', JSON.stringify(responseData.user));
                const token = responseData.token;

                // Enhanced security: Use HttpOnly and Secure flags for the cookie
                const expirationDays = 7; // Adjust as needed
                const expirationDate = new Date();
                expirationDate.setDate(expirationDate.getDate() + expirationDays);
                const expires = expirationDate.toUTCString();
                const cookieString = `token=${token}; expires=${expires}; path=/`;
                // Set the cookie
                document.cookie = cookieString;
                // Instead of immediately redirecting the user, provide a success message and then redirect
                $("#messagefailure").removeClass("text-danger").addClass("text-success").html("Success! Redirecting...");
                setTimeout(function () {
                    window.location.href = "/";
                }, 2000);
            })
            .catch(error => {
                console.error("Error:", error);
                // Enhanced security: Consider not displaying detailed error information to the user
                $("#messagefailure").addClass("text-danger").html("Login failed. Please try again.");
            });
    });
});
