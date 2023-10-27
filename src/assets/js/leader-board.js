const baseUrl = "https://auth.fxchampionship.com";
const urlLeaderBoard = baseUrl + "/auth/contest/get-leaderboard-by-contestid";

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
    const dataLGreeting = JSON.parse(localStorage.getItem("data"))
    let contest_list = dataLGreeting.contest_info_list
    let listContestRender = ""
    let contest_id = ""
    for (let key in contest_list) {
        let class_active = ""
        if (key == 0) {
            class_active = "active";
            contest_id = contest_list[key].contest_id;
        }
        listContestRender += `
            <button class="nav-link ${class_active} fs-4 fw-bolder" onclick="getLeaderBoard('${contest_list[key].contest_id}')" id="${contest_list[key].contest_id}" data-bs-toggle="tab" data-bs-target="#${contest_list[key].contest_id}}"
              type="button" role="tab" aria-controls="${contest_list[key].contest_id}" aria-selected="true">${contest_list[key].contest_id}</button>
        `
    }
    const jwtToken = getCookie("token");
    const inpContest = {
        "contest_id": contest_id
    };

    const headers = new Headers({
        'Authorization': `Bearer ${jwtToken}`
    });
    fetch(urlLeaderBoard, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(inpContest),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(dataResponse => {
            console.log(dataResponse)
            $("#leader_board_title").text(`Leader Board of competition: ${contest_id}`)
            let myRank = dataResponse.rank
            let myData = dataResponse.user
            let curLeaderBoard = dataResponse.leader_board
            // console.log(curLeaderBoard)
            let htmlRender = `
                            <tr class="table-secondary fs-3">
                            <td class="border-bottom-0">
                                <span class="fw-semibold" style="color:#8957FF !important;">${myRank}</span> <br>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-semibold" style="color:#8957FF !important;">${contest_id}</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-semibold" style="color:#8957FF !important;">${myData.login}</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal" style="color:#8957FF !important;">${maskEmail(myData.email)}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal" style="color:#8957FF !important;">$${myData.balance}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0" style="color:#8957FF !important;">$${myData.equity}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0" style="color:#8957FF !important;">$${myData.floating}</span>
                            </td>
                            </tr>
                        `
            for (let key in curLeaderBoard) {
                if (Number(key) > 9) continue;
                htmlRender += `
                            <tr class="fs-3">
                            <td class="border-bottom-0">
                                <span class="fw-normal">${Number(key) + 1}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">${curLeaderBoard[key].contest_id}</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-semibold">${curLeaderBoard[key].login}</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">${maskEmail(curLeaderBoard[key].email)}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">$${curLeaderBoard[key].balance}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0">$${curLeaderBoard[key].equity}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0">$${curLeaderBoard[key].floating}</span>
                            </td>
                            </tr>
                            `
            }
            htmlRender += `
                        <tr class="">
                        <td class="border-bottom-0">
                            <span class="fw-normal">...</span>
                        </td>
                        <td class="border-bottom-0">
                            <span class="fw-normal">...</span>
                        </td>
                        <td class="border-bottom-0">
                            <span class="fw-normal">...</span>
                        </td>
                        <td class="border-bottom-0">
                            <span class="fw-normal">...</span>
                        </td>
                        <td class="border-bottom-0">
                            <span class="fw-normal mb-0">...</span>
                        </td>
                        <td class="border-bottom-0">
                            <span class="fw-normal mb-0">...</span>
                        </td>
                        </tr>
                    `
            $("#leader_board_infor").html(htmlRender)
        })
        .catch(error => {
            console.error("Error:", error);
        });
    $("#nav-tab-contest").html(listContestRender)
})



function splitStringByIndex(inputString, index) {
    if (index >= 0 && index < inputString.length) {
        const firstPart = inputString.substring(0, index); // Get characters from the beginning up to the index
        const secondPart = inputString.substring(index); // Get characters from the index to the end

        return [firstPart, secondPart];
    } else {
        // Index is out of bounds, return an error message or handle the situation as needed.
        return "Invalid index.";
    }
}

function maskEmail(email) {
    let parts = email.split('@');
    if (parts.length === 2) {
        const [firstPart, secondPart] = splitStringByIndex(parts[0], 3);
        let maskedName = firstPart + secondPart.replace(/./g, '*');
        return maskedName + '@' + parts[1];
    } else {
        return email;
    }
}

function getLeaderBoard(contest_id) {
    const jwtToken = getCookie("token");
    const inpContest = {
        "contest_id": contest_id
    };

    const headers = new Headers({
        'Authorization': `Bearer ${jwtToken}`
    });
    fetch(urlLeaderBoard, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(inpContest),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(dataResponse => {
            $("#leader_board_title").text(`Leader Board of competition: ${contest_id}`)
            let myRank = dataResponse.rank
            let myData = dataResponse.user
            let curLeaderBoard = dataResponse.leader_board
            let htmlRender = `
                            <tr class="table-secondary fs-3">
                                <td class="border-bottom-0">
                                    <span class="fw-semibold" style="color:#8957FF !important;">${myRank}</span> <br>
                                </td>
                                <td class="border-bottom-0">
                                    <span class="fw-semibold" style="color:#8957FF !important;">${contest_id}</span> 
                                </td>
                                <td class="border-bottom-0">
                                    <span class="fw-semibold" style="color:#8957FF !important;">${myData.login}</span> 
                                </td>
                                <td class="border-bottom-0">
                                    <span class="fw-normal" style="color:#8957FF !important;">${maskEmail(myData.email)}</span>
                                </td>
                                <td class="border-bottom-0">
                                    <span class="fw-normal" style="color:#8957FF !important;">$${myData.balance}</span>
                                </td>
                                <td class="border-bottom-0">
                                    <span class="fw-normal mb-0" style="color:#8957FF !important;">$${myData.equity}</span>
                                </td>
                                <td class="border-bottom-0">
                                    <span class="fw-normal mb-0" style="color:#8957FF !important;">$${myData.floating}</span>
                                </td>
                            </tr>
                        `
            for (let key in curLeaderBoard) {
                if (Number(key) > 9) continue;
                htmlRender += `
                            <tr class="fs-3">
                            <td class="border-bottom-0">
                                <span class="fw-semibold">${Number(key) + 1}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-semibold">${curLeaderBoard[key].contest_id}</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-semibold">${curLeaderBoard[key].login}</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">${maskEmail(curLeaderBoard[key].email)}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">$${curLeaderBoard[key].balance}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0">$${curLeaderBoard[key].equity}</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0">$${curLeaderBoard[key].floating}</span>
                            </td>
                            </tr>
                            `
            }
            htmlRender += `
                            <tr class="">
                            <td class="border-bottom-0">
                                <span class="fw-semibold">...</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-semibold">...</span> 
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">...</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal">...</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0">...</span>
                            </td>
                            <td class="border-bottom-0">
                                <span class="fw-normal mb-0">...</span>
                            </td>
                            </tr>
                        `
            $("#leader_board_infor").html(htmlRender)
        })
        .catch(error => {
            console.error("Error:", error);
        });
}