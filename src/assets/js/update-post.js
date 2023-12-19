
const baseUrl = "https://admin.fxchampionship.com";
// const baseUrl = "http://localhost:8081";
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
                        ${new Date(datas[key].UpdatedAt).toLocaleString()}
                        </span> 
                      </td>
                      <td class="border-bottom-0">
                        <h6 class="fw-bold p-0 m-0">${datas[key].title}</h6>
                      </td>
                      <td class="border-bottom-0">
                        <img style="height:100px" class="img-inreview" src="${datas[key].thumb}" alt="">
                      </td>
                      <td class="border-bottom-0">
                        <span class="fw-normal">${datas[key].type}</span> <br>
                        <span class="fw-normal"><a target="_blank" href="${'https://fxchampionship.com/' + datas[key].url}">Link</a></span>
                      </td >
                        <td class="border-bottom-0">
                            <span class="fw-normal">${datas[key].viewer}</span>
                        </td>
                        <td class="border-bottom-0">
                            <button onclick="UpdatePost('${datas[key].ID}')" type="button" class="btn btn-primary p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_updates">Update</button>
                        </td>
                        <td class="border-bottom-0">
                            <button onclick="DeletePost('${datas[key].ID}')" type="button" class="btn btn-danger p-1 w-100" data-bs-toggle="modal" data-bs-target="#modal_delete">Delete</button>
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
        plugins: 'anchor autolink charmap codesample emoticons image link lists media searchreplace table visualblocks wordcount',
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
    // Chuyển đổi từng ký tự trong tiêu đề
    var latinTitle = title.split('').map(function (char) {
        return vietnameseToLatin(char);
    }).join('');

    // Chuyển thành chữ thường và giữ lại ký tự chữ cái và số
    var cleanUrl = latinTitle.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/-$/, "");

    return cleanUrl;
}

function vietnameseToLatin(char) {
    // Tạo bảng chuyển đổi giữa tiếng Việt và Latin
    var conversionTable = {
        'à': 'a', 'á': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a',
        'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
        'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
        'è': 'e', 'é': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e',
        'ê': 'e', 'ề': 'e', 'ế': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
        'ì': 'i', 'í': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
        'ò': 'o', 'ó': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o',
        'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
        'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
        'ù': 'u', 'ú': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u',
        'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
        'ỳ': 'y', 'ý': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
        'đ': 'd',
        ' ': '-' // Chuyển khoảng trắng thành dấu gạch ngang
    };

    // Trả về ký tự Latin tương ứng hoặc ký tự gốc nếu không có trong bảng
    return conversionTable[char] || char;
}

function UpdatePost(id_param) {
    let url = baseUrl + "/public/post-by-id?id=" + id_param
    // console.log(url)
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
            $("#post_to_server").prop("disabled", false);
            $("#fb_msg_create").text("")
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
                $("#post_to_server").prop("disabled", true);
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
                // console.log(JSON.stringify(post_data))
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
                        $("#fb_msg_create").addClass("text-success").text("Success!")
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

function DeletePost(id_param) {
    $("#delete_to_server").prop("disabled", false);
    $("#fb_msg_delete").addClass("text-success").text("")
    let url = baseUrl + "/public/post-by-id?id=" + id_param
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
                $("#delete_to_server").prop("disabled", true);
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
                        $("#fb_msg_delete").addClass("text-success").text("Success!")
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