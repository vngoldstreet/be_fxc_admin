
const urlCompetition = "https://auth.fxchampionship.com/auth/contest/get-contest-by-uid";
const urlRejoinCompetition = "https://auth.fxchampionship.com/auth/contest/rejoin-a-competition";

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
    const jwtToken = getCookie("token");
    const headers = new Headers({
        'Authorization': `Bearer ${jwtToken}`
    });
    fetch(urlCompetition, {
        method: "GET",
        headers: headers,
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(dataResponse => {
            // console.log(dataResponse)
            let myData = dataResponse.data
            let htmlRender = ``
            if (myData.length == 0) {
                $("#competition").html("<p class='text-danger'>You haven't participated in any competition yet.</p>");
                return
            }

            for (let key in myData) {
                htmlRender += `
                <div class="col-lg-3 col-md-6 mt-3">
                    <img class="w-100 rounded" src="src/assets/images/fxcentrum-live-Contest.jpg">
                    <h5 class="fw-bolder mt-3 text-violet-fxc">Competition ID: ${myData[key].contest_id}</h5>
                    <p class="card-text"><span class="fw-bolder">FXID:</span> ${myData[key].fx_id}</p>
                    <p class="card-text"><span class="fw-bolder">Server:</span>FXU Solution</p>
                    <p class="card-text" type="password"><span class="fw-bolder">Master Password:</span> <input id="master-${myData[key].contest_id}" type="password" class="border-0" value="${myData[key].fx_master_pw}"></p>
                    <p class="card-text" type="password"><span class="fw-bolder">Master Password:</span> <input id="invester-${myData[key].contest_id}" type="password" class="border-0" value="${myData[key].fx_invester_pw}"></p>
                    <div class="d-flex justify-content-between">
                    <button id="btn-${myData[key].contest_id}" type="button" onclick="showpassword('${myData[key].contest_id}')" class="btn btn-outline-primary">Show password</button>
                    <button type="button" onclick = "confirmToReJoin('${myData[key].contest_id}')" class="btn btn-outline-success" data-bs-toggle="modal" data-bs-target="#re_join_contest">Join again</button>
                    </div>
                </div>
            `
            }
            $("#competition").html(htmlRender);
        })
        .catch(error => {
            console.error("Error:", error);
        });
})

// onclick = "reJoinACompetition('${myData[key].contest_id}')"

function showpassword(param) {
    let id_master = "#master-" + param
    let id_invester = "#invester-" + param
    let btn_id = "#btn-" + param
    var passwordFieldType = $(id_master).attr('type');
    if (passwordFieldType === 'password') {
        $(id_master).attr('type', 'text');
        $(id_invester).attr('type', 'text');
        $(btn_id).text("Hide password")
        setTimeout(function () {
            $(id_master).attr('type', 'password');
            $(id_invester).attr('type', 'password');
            $(btn_id).text("Show Password")
        }, 15000)
    } else {
        $(id_master).attr('type', 'password');
        $(id_invester).attr('type', 'password');
        $(btn_id).text("Show Password")
    }
}

function confirmToReJoin(contest_id) {
    $("#this_contest_info").remove();
    let html_text = `
        <h6><span class="fw-semibold">ID:</span> ${contest_id}</h6>
    `
    $("#contest_info").html(html_text);
    $("#confirm_to_re_join").click(function () {
        reJoinACompetition(contest_id);
    });
}

function reJoinACompetition(contest_id) {
    const jwtToken = getCookie("token");
    const headers = new Headers({
        'Authorization': `Bearer ${jwtToken}`
    });
    const rejoinContest = {
        "contest_id": contest_id
    };
    console.log(urlRejoinCompetition + ": " + JSON.stringify(rejoinContest))
    fetch(urlRejoinCompetition, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(rejoinContest),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(dataResponse => {
            window.alert("Success!")
        })
        .catch(error => {
            window.alert(error.message)
            console.error("Error:", error);
        });
}