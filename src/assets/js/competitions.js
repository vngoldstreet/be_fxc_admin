let urlCompetitionList = "/auth/get-competition-request-list";
let urlCreateALoginID = "/auth/contest-approval";
let urlRejoinToContest = "/auth/rejoin-contest-approval";
let urlConfirmationTransactions = "/auth/admin-transaction";
let urlRejectTransactions = "/auth/cancel-transaction";
let goldRate = 24000;

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
  let jwtToken = getCookie("token");

  if (!jwtToken) {
    redirectToURL('/login');
    return;
  }

  let headers = new Headers({
    "Content-Type": "application/json",
    'Authorization': `Bearer ${jwtToken}`
  });

  fetch(urlCompetitionList, {
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
      let transactionData = dataResponse.data
      for (let key in transactionData) {
        let text_type = "";
        let text_id_contest = "";
        let btn_confirm = ""

        switch (transactionData[key].type_id) {
          case 1:
            text_type = "Deposit";
            break;
          case 2:
            text_type = "Withdrawal";
            break;
          case 3:
            text_type = "Promo";
            break;
          case 4:
            text_type = "Join a contest";
            text_id_contest = `${transactionData[key].contest_id}`;
            btn_confirm = `<button onclick="ShowTransactionInformation(${transactionData[key].ID},${transactionData[key].customer_id},'${text_id_contest}','${transactionData[key].name}')" type="button" class="btn btn-secondary p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_transaction_info">Accept to Join</button>`
            break;
          case 5:
            text_type = "Earning";
            break;
          case 6:
            text_type = "Re-Join a contest";
            text_id_contest = `${transactionData[key].contest_id}`;
            btn_confirm = `<button onclick="ShowRejoin(${transactionData[key].ID},${transactionData[key].customer_id},'${text_id_contest}','${transactionData[key].name}')" type="button" class="btn btn-outline-success p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_transaction_rejoin">Accept to Re-Join</button>`
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
        let updated_at = new Date(transactionData[key].UpdatedAt).toLocaleString();
        let created_at = new Date(transactionData[key].CreatedAt).toLocaleString();
        let number = Number(key) + 1;
        let amount = Number(transactionData[key].amount).toLocaleString();
        let vndAmount = Number(transactionData[key].amount) * goldRate

        htmlPrint += `
                    <tr>
                      <td class="border-bottom-0">
                        <span class="fw-semibold">${number}/${transactionData[key].ID}</span>
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
                      <td class="border-bottom-0">
                        ${btn_confirm}
                      </td>
                      <td class="border-bottom-0">
                        <button onclick="CancelTransaction(${transactionData[key].ID},'${transactionData[key].name}','${amount}','${text_id_contest}')" type="button" class="btn btn-danger p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_transaction_reject">Reject</button>
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

function ShowTransactionInformation(transaction_id, customer_id, contest_id, name) {
  $("#confirm_for_transaction").prop("disabled", false);
  $("#confirm_for_transaction").text(`Approval for this competition: ${contest_id}`)
  $("#inpContestID").attr("value", contest_id)
  $("#inpCustomerID").attr("value", `${customer_id} - ${name}`)

  $("#confirm_for_transaction").on("click", function (e) {
    e.preventDefault()
    $("#confirm_for_transaction").prop("disabled", true);

    // let fx_id_text = $("#inpLoginID").val();
    // let fx_master_password = $("#inpMasterPassword").val();
    // let fx_investor_password = $("#inpInvestorPassword").val();

    // if (fx_id_text === '') {
    //   $('#inpLoginID').addClass('is-invalid');
    //   $('#fb_fx_id_text').addClass('invalid-feedback').text('LoginID is required'); // Display an error message
    //   return;
    // } else {
    //   // Valid email format
    //   $('#inpLoginID').removeClass('is-invalid').addClass('is-valid');
    //   $('#fb_fx_id_text').removeClass('invalid-feedback').addClass('invalid-feedback').text('Look good'); // Clear the error message
    // }
    // if (fx_master_password === '') {
    //   $('#inpMasterPassword').addClass('is-invalid');
    //   $('#fb_fx_master_password').addClass('invalid-feedback').text('MasterPassword is required'); // Display an error message
    //   return;
    // } else {
    //   // Valid email format
    //   $('#inpMasterPassword').removeClass('is-invalid').addClass('is-valid');
    //   $('#fb_fx_master_password').removeClass('invalid-feedback').addClass('invalid-feedback').text('Look good'); // Clear the error message
    // }
    // if (fx_investor_password === '') {
    //   $('#inpInvestorPassword').addClass('is-invalid');
    //   $('#fb_fx_invester_password').addClass('invalid-feedback').text('InvestorPassword is required'); // Display an error message
    //   return;
    // } else {
    //   // Valid email format
    //   $('#inpInvestorPassword').removeClass('is-invalid').addClass('is-valid');
    //   $('#fb_fx_invester_password').removeClass('invalid-feedback').addClass('invalid-feedback').text('Look good'); // Clear the error message
    // }

    //, fx_id_text, fx_master_password, fx_investor_password
    CreateMetaTraderData(contest_id, customer_id, transaction_id)
  });
}

function ShowRejoin(param_id, customer_id, contest_id, name) {
  $("#approval_title_rejoin").text(`Approval for this competition: ${contest_id}`)
  $("#inpContestIDRejoin").attr("value", contest_id)
  $("#inpCustomerIDRejoin").attr("value", `${customer_id} - ${name}`)

  $("#confirm_for_rejoin").on("click", function (e) {
    e.preventDefault()
    ApprovalRejoinToContest(contest_id, customer_id)
    ConfirmTransaction(param_id);
  });
}

function ApprovalRejoinToContest(contest_id, customer_id) {
  let jwtToken = getCookie("token");
  if (!jwtToken) {
    console.error("Error: JWT token is missing.");
    return;
  }

  let inpApproval = {
    "contest_id": contest_id,
    "customer_id": customer_id
  };

  let headers = new Headers({
    'Authorization': `Bearer ${jwtToken}`
  });
  // console.log(JSON.stringify(inpApproval))
  fetch(urlRejoinToContest, {
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
      $("#re_join_contest_message").removeClass().addClass("fw-semibold");
      $("#re_join_contest_message").html(html);
    })
    .catch(error => {
      console.error("Error:", error);
    });
}
//, fx_id, fx_master_password, fx_investor_password
function CreateMetaTraderData(contest_id, customer_id, transaction_id) {
  let jwtToken = getCookie("token");
  if (!jwtToken) {
    console.error("Error: JWT token is missing.");
    return;
  }

  let inpApproval = {
    "contest_id": contest_id,
    "customer_id": customer_id
    // "fx_id": fx_id,
    // "fx_master_pw": fx_master_password,
    // "fx_invester_pw": fx_investor_password,
  };
  // console.log(JSON.stringify(inpApproval))
  let headers = new Headers({
    'Authorization': `Bearer ${jwtToken}`
  });

  fetch(urlCreateALoginID, {
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
      $("#join_contest_message").removeClass().addClass("fw-semibold");
      $("#join_contest_message").html(html);
      // ConfirmTransaction(transaction_id);
      GetListOfTransactions()
    })
    .catch(error => {
      let html = `<code class='w-100 text-danger'>${error}</code>`
      $("#join_contest_message").html(html);
      console.error("Error:", error);
    });
}

function ConfirmRejectTransaction(param_id) {
  let jwtToken = getCookie("token");
  if (!jwtToken) {
    console.error("Error: JWT token is missing.");
    return;
  }

  let inpReject = {
    "id": param_id
  };

  let headers = new Headers({
    'Authorization': `Bearer ${jwtToken}`
  });

  fetch(urlRejectTransactions, {
    method: "POST",
    headers: headers,
    body: JSON.stringify(inpReject),
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
      $("#reject-transaction").removeClass().addClass("fw-semibold");
      $("#reject-transaction").html(html);
      GetListOfTransactions()
    })
    .catch(error => {
      console.error("Error:", error);
    });
}


function CancelTransaction(param_id, name, amount, contest_id) {
  $("#reject_for_transaction").prop("disabled", false);
  let vndAmount = amount * goldRate
  let html_text = `
  <p><span class="fw-semibold">ID:</span> ${param_id}</p>
  <p><span class="fw-semibold">ContestID:</span> ${contest_id}</p>
  <p><span class="fw-semibold">Name:</span> ${name}</p>
  <p><span class="fw-semibold">Amount:</span> ${amount.toLocaleString()} G - ${vndAmount.toLocaleString()} VND</p>
  <p id="reject-transaction" class="fw-semibold"></p>
  `

  $("#transaction-information-reject").html(html_text)
  $("#reject_for_transaction").on("click", function (e) {
    e.preventDefault()
    $("#reject_for_transaction").prop("disabled", true);
    ConfirmRejectTransaction(param_id);
  });
}


function ConfirmTransaction(param_id) {
  let jwtToken = getCookie("token");
  if (!jwtToken) {
    console.error("Error: JWT token is missing.");
    return;
  }

  let inpApproval = {
    "id": param_id
  };

  let headers = new Headers({
    'Authorization': `Bearer ${jwtToken}`
  });
  // console.log(JSON.stringify(inpApproval))
  fetch(urlConfirmationTransactions, {
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
      $("#join_contest_message").removeClass().addClass("fw-semibold");
      $("#join_contest_message").html(html);
      GetListOfTransactions()
    })
    .catch(error => {
      console.error("Error:", error);
    });
}