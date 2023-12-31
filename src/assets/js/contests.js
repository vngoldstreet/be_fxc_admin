let urlContestList = "/auth/get-contest-list";
let urlUpdateContest = "/auth/update-contest-id";
let urlCreatedContest = "/auth/create-contest";

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
  $("#create_contest").prop("disabled", false);
  let jwtToken = getCookie("token");

  if (!jwtToken) {
    redirectToURL('/login');
    return;
  }

  let headers = new Headers({
    "Content-Type": "application/json",
    'Authorization': `Bearer ${jwtToken}`
  });

  fetch(urlContestList, {
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
      // console.log(contestDatas)
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
        let type_contest = ""
        switch (contestDatas[key].type_id) {
          case 1:
            type_contest = "Silver";
            break;
          case 2:
            type_contest = "Gold";
            break;
          case 3:
            type_contest = "Platinum";
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
                        <span class="fw-normal mb-0">${contestDatas[key].amount} G</span> <br>
                        <span class="fw-normal mb-0">${type_contest}</span> <br>
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
                        <button onclick="EditingContest('${contestDatas[key].contest_id}','${contestDatas[key].status_id}','${contestDatas[key].type_id}')" type="button" class="btn btn-danger p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_editing">Editing</button>
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

  $("#create_contest").on("click", function (e) {
    e.preventDefault()
    $("#create_contest").prop("disabled", true);
    let amount = $("#inpAmount").val()
    let max_person = $("#inpMaximumPerson").val()
    let start_balance = $("#inpStartBalance").val()
    let start_at = $("#inpStartAt").val()
    let expired_at = $("#inpExpireAt").val()
    let status_id = $("#inpCreateStatusID").val()
    let type_id = $("#inpCreatTypeID").val()

    var datestart = new Date(start_at);
    datestart.setHours(datestart.getHours() + 7);
    var extractedDateStart = datestart.toISOString().split('T')[0];
    var extractedTimeStart = datestart.toISOString().split('T')[1].split('Z')[0];
    var time_start_at = extractedDateStart + " " + extractedTimeStart

    var dateend = new Date(expired_at);
    dateend.setHours(dateend.getHours() + 7);
    var extractedDateEnd = dateend.toISOString().split('T')[0];
    var extractedTimeEnd = dateend.toISOString().split('T')[1].split('Z')[0];
    var time_end = extractedDateEnd + " " + extractedTimeEnd

    let inpCreate = {
      "amount": Number(amount),
      "maximum_person": Number(max_person),
      "start_balance": Number(start_balance),
      "start_at": time_start_at,
      "expired_at": time_end,
      "status_id": Number(status_id),
      "type_id": Number(type_id),
    };

    console.log(JSON.stringify(inpCreate))

    let jwtToken = getCookie("token");
    if (!jwtToken) {
      console.error("Error: JWT token is missing.");
      return;
    }

    let headers = new Headers({
      'Authorization': `Bearer ${jwtToken}`
    });

    fetch(urlCreatedContest, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(inpCreate),
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
        GetListOfContests()
      })
      .catch(error => {
        let stringData = JSON.stringify(dataResponse)
        let html = `<code class='w-100 text-danger'>${stringData}</code>`
        $("#fb_msg_create").removeClass().addClass("fw-semibold");
        $("#fb_msg_create").html(html);
        console.error("Error:", error);
      });
  })
})


function EditingContest(contest_id, status_id, type_id) {
  $("#confirm_for_contest").prop("disabled", false);
  $("#approval_title").text(`Update for this competition: ${contest_id}`)
  $("#inpContestID").attr("value", contest_id)
  $("#inpStatusID").val(status_id)
  $("#inpTypeID").val(type_id)

  $("#confirm_for_contest").on("click", function (e) {
    e.preventDefault()
    $("#confirm_for_contest").prop("disabled", true);
    let stID = $("#inpStatusID").val();
    let newInpTypeID = $("#inpTypeID").val();
    let jwtToken = getCookie("token");
    if (!jwtToken) {
      console.error("Error: JWT token is missing.");
      return;
    }

    let inpApproval = {
      "contest_id": contest_id,
      "status_id": Number(stID),
      "type_id": Number(newInpTypeID)
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


