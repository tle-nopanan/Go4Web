const baseURL = "http://localhost:5000/course";
fetch(baseURL)
  .then((resp) => {
    return resp.json();
  })
  .then((data) => appendData(data))
  .catch((err) => {
    console.log(err);
  });

function appendData(data) {
  var mainContainer = document.getElementById("myData");
  for (var i = 0; i < data.length; i++) {
    var div = document.createElement("div");
    div.innerHTML =
      "CourseID: " +
      data[i].ID +
      " " +
      data[i].Name +
      " " +
      data[i].Price +
      " " +
      data[i].Instructor;
    mainContainer.appendChild(div);
  }
}
