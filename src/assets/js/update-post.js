
// const baseUrl = "https://admin.fxchampionship.com";
const baseUrl = "http://localhost:8081";
const urlUpdate = baseUrl + "/auth/update-post";
const urlDelete = baseUrl + "/auth/delete-post?id=";
const urlGetListOfPost = baseUrl + "/public/posts";

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
    GetStarted();
});

function GetStarted() {
    let url = urlGetListOfPost + "?type=news&length=" + 100
    fetch(url, {
        method: "GET",
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not successful");
            }
            return response.json(); // Parse the response JSON if needed
        })
        .then(dataResponse => {
            console.log(dataResponse)
            $("#total_count").text(`List of ${dataResponse.count} posts`)
            let htmlPrint = "";
            let datas = dataResponse.data
            for (let key in datas) {
                htmlPrint += `
                    <tr>
                      <td class="border-bottom-0">
                        <span class="fw-semibold">${datas[key].ID}</span> <br>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-semibold">
                        ${new Date(datas[key].CreatedAt).toLocaleString()}<br>
                        ${new Date(datas[key].UpdatedAt).toLocaleString()}
                        </span> 
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-bold">${datas[key].title}</span> <br>
                        <span class="fw-normal">${datas[key].description}</span>
                      </td>
                      <td class="border-bottom-0">
                        <img class="img-inreview" src="${baseUrl + '/public/image/' + datas[key].url}" alt="">
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${datas[key].keyword}</span> <br>
                        <span class="fw-normal">${datas[key].tag}</span>
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${datas[key].type}</span> <br>
                        <span class="fw-normal">${baseUrl + '/' + datas[key].url}</span>
                      </td >
                        <td class="border-bottom-0">
                            <span class="fw-normal">${datas[key].viewer}</span>
                        </td>
                        <td class="border-bottom-0">
                            <button onclick="UpdatePost('${datas[key].url}')" type="button" class="btn btn-primary p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_updates">Update</button>
                        </td>
                        <td class="border-bottom-0">
                            <button onclick="DeletePost('${datas[key].url}')" type="button" class="btn btn-danger p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_delete">Delete</button>
                        </td>
                    </tr >
                    `;
            }
            $("#transaction-list").html(htmlPrint);
        })
        .catch(error => {
            console.error("Error:", error);
        });
    tinymce.init({
        selector: '#input_content',
        plugins: 'ai tinycomments mentions anchor autolink charmap codesample emoticons image link lists media searchreplace table visualblocks wordcount checklist mediaembed casechange export formatpainter pageembed permanentpen footnotes advtemplate advtable advcode editimage tableofcontents mergetags powerpaste tinymcespellchecker autocorrect a11ychecker typography inlinecss',
        toolbar: 'undo redo | blocks fontfamily fontsize | bold italic underline strikethrough | link image media table mergetags | align lineheight | tinycomments | checklist numlist bullist indent outdent | emoticons charmap | removeformat',
        tinycomments_mode: 'embedded',
        tinycomments_author: 'Author name',
        mergetags_list: [
            { value: 'First.Name', title: 'First Name' },
            { value: 'Email', title: 'Email' },
        ],
        ai_request: (request, respondWith) => respondWith.string(() => Promise.reject("See docs to implement AI Assistant")),
    });
}

function cleanAndGenerateUrl(title) {
    var cleanUrl = title.toLowerCase().replace(/[^a-z0-9]+/g, '-');
    return cleanUrl;
}

function UpdatePost(url_param) {
    let url = baseUrl + "/public/post-by-url?url=" + url_param
    fetch(url, {
        method: "GET",
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not successful");
            }
            return response.json(); // Parse the response JSON if needed
        })
        .then(dataResponse => {
            let saveData = dataResponse.data
            tinymce.activeEditor.setContent(saveData.content)
            $("#input_title").val(saveData.title);
            $("#input_description").val(saveData.description);
            $("#input_tag").val(saveData.tag);
            $("#input_view").val(saveData.viewer);
            $("#input_type").val(saveData.type);
            $("#input_keyword").val(saveData.keyword);
            $("#input_url").val(saveData.url);
            $("#previewImage").attr("src", saveData.thumb);

            let output_thumb = ""
            $("#input_thumb").change(function () {
                var fileInput = this;
                if (fileInput.files && fileInput.files[0]) {
                    var reader = new FileReader();
                    reader.onload = function (e) {
                        $("#previewImage").attr("src", e.target.result);
                        output_thumb = e.target.result
                    };
                    reader.readAsDataURL(fileInput.files[0]);
                }
            });

            let output_url = ""
            $("#generate_url").click(function () {
                var title = $("#input_title").val();
                var resultUrl = cleanAndGenerateUrl(title);
                $("#input_url").val(resultUrl);
                output_url = resultUrl
            });

            $("#post_to_server").click(function () {
                let output_title = $("#input_title").val();
                let output_description = $("#input_description").val();
                let output_tag = $("#input_tag").val();
                let output_view = $("#input_view").val();
                let output_type = $("#input_type").val();
                let output_content = tinymce.activeEditor.getContent();
                let output_keyword = $("#input_keyword").val();
                if (output_url === "") {
                    output_url = cleanAndGenerateUrl(output_title);
                }
                let post_data = {
                    "id": saveData.ID,
                    "title": output_title,
                    "description": output_description,
                    "content": output_content,
                    "type": output_type,
                    "thumb": output_thumb,
                    "tag": output_tag,
                    "viewer": Number(output_view),
                    "url": output_url,
                    "keyword": output_keyword
                };

                const jwtToken = getCookie("token");
                if (!jwtToken) {
                    console.error("Error: JWT token is missing.");
                    return;
                }

                const headers = new Headers({
                    'Authorization': `Bearer ${jwtToken}`
                });
                console.log(JSON.stringify(post_data))
                fetch(urlUpdate, {
                    method: "POST",
                    headers: headers,
                    body: JSON.stringify(post_data),
                })
                    .then(response => {
                        if (!response.ok) {
                            throw new Error("Network response was not successful");
                        }
                        return response.json(); // Parse the response JSON if needed
                    })
                    .then(dataResponse => {
                        $("#fb_msg_create").addClass("text-success").text(JSON.stringify(dataResponse.data))
                        GetStarted()
                    })
                    .catch(error => {
                        console.error("Error:", error);
                    });

            })
        })
        .catch(error => {
            console.error("Error:", error);
        });
}

function DeletePost(url_param) {
    let url = baseUrl + "/public/post-by-url?url=" + url_param
    fetch(url, {
        method: "GET",
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not successful");
            }
            return response.json(); // Parse the response JSON if needed
        })
        .then(dataResponse => {
            let saveData = dataResponse.data
            $("#title_delete").text(saveData.title);
            $("#delete_desc").text(saveData.description);
            $("#delete_thumb").html(`<img class="img-inreview" src="${baseUrl + '/public/image/' + saveData.url}" alt="">`);
            $("#delete_to_server").click(function () {
                const jwtToken = getCookie("token");
                if (!jwtToken) {
                    console.error("Error: JWT token is missing.");
                    return;
                }
                const headers = new Headers({
                    'Authorization': `Bearer ${jwtToken}`
                });
                let url_delete = urlDelete + saveData.ID
                fetch(url_delete, {
                    method: "POST",
                    headers: headers,
                })
                    .then(response => {
                        if (!response.ok) {
                            throw new Error("Network response was not successful");
                        }
                        return response.json(); // Parse the response JSON if 
                    })
                    .then(dataResponse => {
                        $("#fb_msg_delete").addClass("text-success").text(JSON.stringify(dataResponse))
                        GetStarted()
                    })
                    .catch(error => {
                        console.error("Error:", error);
                    });

            })
        })
        .catch(error => {
            console.error("Error:", error);
        });
}