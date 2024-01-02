let urlHistoryContestList = "/auth/get-history-contest-list";
let urlUpdateContest = "/auth/update-contest-id";

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

function GetListOfContests() {
    let jwtToken = getCookie("token");

    if (!jwtToken) {
        redirectToURL('/login');
        return;
    }

    let headers = new Headers({
        "Content-Type": "application/json",
        'Authorization': `Bearer ${jwtToken}`
    });

    fetch(urlHistoryContestList, {
        method: "GET",
        headers: headers
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json(); // Parse the response JSON if needed
        })
        .then(dataResponse => {
            let htmlPrint = "";
            let contestDatas = dataResponse.data
            let text_status = "";
            let bg_class = "";
            for (let key in contestDatas) {
                switch (contestDatas[key].status_id) {
                    case 0:
                        text_status = "Future";
                        bg_class = "text-primary";
                        break;
                    case 1:
                        text_status = "Processing";
                        bg_class = "text-warning";
                        break;
                    case 2:
                        text_status = "Finished";
                        bg_class = "text-success";
                        break;
                    case 3:
                        text_status = "Cancel";
                        bg_class = "text-danger";
                        break;
                    default:
                        break;
                }
                htmlPrint += `
                    <tr>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${contestDatas[key].ID}</span><br>
                        <span class="fw-normal">${contestDatas[key].contest_id}</span>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${new Date(contestDatas[key].CreatedAt).toLocaleString()}</span> <br>
                        <span class="fw-normal">${new Date(contestDatas[key].UpdatedAt).toLocaleString()}</span>
                      </td>r
                      <td class="border-bottom-0">
                        <span class="fw-normal">${new Date(contestDatas[key].start_at).toLocaleString()}</span> <br>
                        <span class="fw-normal">${new Date(contestDatas[key].expired_at).toLocaleString()}</span> 
                      </td>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${contestDatas[key].current_person}/</span><br>
                        <span class="fw-normal">${contestDatas[key].maximum_person}</span>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal mb-0">${contestDatas[key].amount} G</span>
                      </td>
                      <td class="border-bottom-0">
                        <div class="d-flex align-items-center gap-2">
                          <span class="fw-normal mb-0">$${contestDatas[key].start_balance}</span>
                        </div>
                      </td>
                      <td class="border-bottom-0">
                        <div class="d-flex align-items-center gap-2">
                          <span class="badge ${bg_class} rounded-1 fw-semibold">${text_status}</span>
                        </div>
                      </td>
                      <td class="border-bottom-0">
                        <button onclick="EditingContest('${contestDatas[key].contest_id}','${contestDatas[key].status_id}')" type="button" class="btn btn-danger p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_editing">Editing</button>
                      </td>
                    </tr>
                  `;
            }

            $("#contests-list").html(htmlPrint);
        })
        .catch(error => {
            console.error("Error:", error);
        });
}

$(document).ready(function () {
    GetListOfContests()
})

function EditingContest(contest_id, status_id) {
    $("#approval_title").text(`Update for this competition: ${contest_id}`)
    $("#inpContestID").attr("value", contest_id)
    $("#inpStatusID").val(status_id)

    $("#confirm_for_contest").click(function () {
        let stID = $("#inpStatusID").val();
        let new_status_id = Number(stID)
        let jwtToken = getCookie("token");
        if (!jwtToken) {
            console.error("Error: JWT token is missing.");
            return;
        }

        let inpApproval = {
            "contest_id": contest_id,
            "status_id": new_status_id,
        };

        let headers = new Headers({
            'Authorization': `Bearer ${jwtToken}`
        });
        console.log(JSON.stringify(inpApproval))
        fetch(urlUpdateContest, {
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
                let stringData = JSON.stringify(dataResponse)
                let html = `<code class='w-100 text-success'>${stringData}</code>`
                $("#fb_msg").removeClass().addClass("fw-semibold");
                $("#fb_msg").html(html);
                GetListOfContests()
            })
            .catch(error => {
                console.error("Error:", error);
            });
    });
}
