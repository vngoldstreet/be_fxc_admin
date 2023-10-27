// const baseUrl = "https://admin.fxchampionship.com";
const baseUrl = "http://localhost:8081";
const urlTransactionList = baseUrl + "/auth/get-history-transaction-list";
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

function GetListOfTransactions() {
    const jwtToken = getCookie("token");

    if (!jwtToken) {
        redirectToURL('/login');
        return;
    }

    const headers = new Headers({
        "Content-Type": "application/json",
        'Authorization': `Bearer ${jwtToken}`
    });

    fetch(urlTransactionList, {
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
            const userInfo = JSON.parse(localStorage.getItem('user'));
            let transactionData = dataResponse.data
            for (let key in transactionData) {
                let text_type = "";
                let text_id_contest = "";
                switch (transactionData[key].type_id) {
                    case 1:
                        text_type = "Deposit";
                        break;
                    case 2:
                        text_type = "Withdrawal";
                        break;
                    case 3:
                    case 5:
                        text_type = "Earning";
                        break;
                    case 4:
                        text_type = "Join a contest";
                        text_id_contest = `${transactionData[key].contest_id}`;
                        break;
                    default:
                        break;
                }
                let text_status = "";
                let text_class = "";
                let bg_class = "";
                switch (transactionData[key].status_id) {
                    case 1:
                        text_status = "Processing";
                        text_class = "text-warning";
                        bg_class = "bg-warning";
                        break;
                    case 2:
                        text_status = "Success";
                        text_class = "text-success";
                        bg_class = "bg-success";
                        break;
                    case 3:
                        text_status = "Cancelled";
                        text_class = "text-danger";
                        bg_class = "bg-danger";
                        break;
                }
                const updated_at = new Date(transactionData[key].UpdatedAt).toLocaleString();
                const created_at = new Date(transactionData[key].CreatedAt).toLocaleString();
                const number = Number(key) + 1;
                const amount = Number(transactionData[key].amount).toLocaleString();
                let vndAmount = Number(transactionData[key].amount) * goldRate
                htmlPrint += `
                    <tr>
                      <td class="border-bottom-0">
                        <span class="fw-semibold">${number}/${transactionData[key].ID}</span> <br>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-semibold">${text_type}<br>${text_id_contest}</span> 
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${created_at}</span> <br>
                        <span class="fw-normal">${updated_at}</span>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${transactionData[key].name}</span> <br>
                        <span class="fw-normal">${transactionData[key].email}</span><br>
                        <span class="fw-normal">${transactionData[key].phone}</span>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal mb-0">${amount} G</span> <br>
                        <span class="fw-normal mb-0">${vndAmount.toLocaleString()} VND</span>
                      </td>
                      <td class="border-bottom-0">
                        <div class="d-flex align-items-center gap-2">
                          <span class="badge ${bg_class} rounded-1 fw-semibold">${text_status}</span>
                        </div>
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
    GetListOfTransactions()
})
