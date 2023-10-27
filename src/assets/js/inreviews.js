// const baseUrl = "https://admin.fxchampionship.com";
const baseUrl = "http://localhost:8081";
const urlInreviewList = baseUrl + "/auth/get-review-list";
const urlInreviewUpdate = baseUrl + "/auth/update-review-list";
const goldRate = 24000;

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

function redirectToURL(targetUrl) {
    window.location.href = targetUrl;
}

function GetListInReviews() {
    const jwtToken = getCookie("token");

    if (!jwtToken) {
        redirectToURL('/login');
        return;
    }

    const headers = new Headers({
        "Content-Type": "application/json",
        'Authorization': `Bearer ${jwtToken}`
    });

    fetch(urlInreviewList, {
        method: "GET",
        headers: headers
    })
        .then(response => {
            if (!response.ok) {
                console.log(response)
                throw new Error("Network response was not ok");
            }
            return response.json(); // Parse the response JSON if needed
        })
        .then(dataResponse => {
            let htmlPrint = "";
            let transactionData = dataResponse.data
            for (let key in transactionData) {
                htmlPrint += `
                    <tr>
                      <td class="border-bottom-0">
                        <span class="fw-semibold">${transactionData[key].ID}</span> <br>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${new Date(transactionData[key].CreatedAt).toLocaleString()}</span> <br>
                        <span class="fw-normal">${new Date(transactionData[key].UpdatedAt).toLocaleString()}</span>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${transactionData[key].name}</span> <br>
                        <span class="fw-normal">${transactionData[key].email}</span><br>
                        <span class="fw-normal">${transactionData[key].phone}</span>
                      </td>
                      <td class="border-bottom-0">
                        <img onclick="ShowImage('${transactionData[key].image_front}')" style="height: 100px;aspect-ratio: 16/9;object-fit: cover;width: 100%;cursor: pointer;" src="${transactionData[key].image_front}" data-bs-toggle="modal" data-bs-target="#modal_image"/>
                      </td>
                       <td class="border-bottom-0">
                        <img  onclick="ShowImage('${transactionData[key].image_front}')" style="height: 100px;aspect-ratio: 16/9;object-fit: cover;width: 100%;cursor: pointer;" class="img-inreview" src="${transactionData[key].image_back}" data-bs-toggle="modal" data-bs-target="#modal_image"/>
                      </td>
                      <td class="border-bottom-0">
                        <div class="d-flex align-items-center gap-2">
                          <span class="badge text-warning rounded-1 fw-semibold">${transactionData[key].status}</span>
                        </div>
                      </td>
                      <td class="border-bottom-0">
                        <button onclick="Confirmation(${transactionData[key].customer_id},'done')" type="button" class="btn btn-danger p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_confirm_image">Confirmation</button>
                      </td>
                    </tr>
                  `;
            }
            $("#transaction-list").html(htmlPrint);
        })
        .catch(error => {
            console.error("Error:", error);
        });
}

$(document).ready(function () {
    GetListInReviews()
})


function ShowImage(param) {
    $("#how_img_title").text(`Show image`)
    $("#img-to-show").attr("src", param)
}

function Confirmation(customer_id, status) {
    $("#inpCustomerID").attr("value", customer_id)
    $("#inpStatus").attr("value", status)
    $("#confirm_for_inreview").click(function () {
        SendRequest(customer_id, status);
    });
}

function SendRequest(customer_id, status) {
    const jwtToken = getCookie("token");
    if (!jwtToken) {
        console.error("Error: JWT token is missing.");
        return;
    }

    const inpApproval = {
        "customer_id": customer_id,
        "status": status,
    };

    const headers = new Headers({
        'Authorization': `Bearer ${jwtToken}`
    });
    console.log(JSON.stringify(inpApproval))
    fetch(urlInreviewUpdate, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(inpApproval),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not successful");
            }
            return response.json(); // Parse the response JSON if needed
        })
        .then(dataResponse => {
            let html = `<code class='w-100 text-success'>Success!</code>`
            $("#msg_inreview").removeClass().addClass("fw-semibold");
            $("#msg_inreview").html(html);
        })
        .catch(error => {
            console.error("Error:", error);
        });
}