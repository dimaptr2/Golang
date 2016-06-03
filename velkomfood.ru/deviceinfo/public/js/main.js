/*
    Global variables
*/

// request parameters
var reqFrom
var reqTo

/*
    Functions
*/

// Initialization of AJAX in the page
// It is job only under Firefox, Opera, Safari and, Chrome or Chromium browsers.

// Set the format of the date.
function getDateFormat() {

    var dt = new Date()

    var day = dt.getDate()
    var month = dt.getMonth()
    var year = dt.getFullYear()
    var strDay = ""
    var strMonth = ""

    if (day > 1) {
        day = day - 1
    }


    if (day >= 1 && day <=9) {
        strDay = "0" + day
    } else {
        strDay = day
    }

    if (month < 12) {
        month = month + 1
    }

    if (month >= 1 && month <= 9) {
        strMonth = "0" + month
    } else {
        strMonth = month
    }

    return (strDay + "." + strMonth + "." + year)


}

// Set values in the input element.

function getDateFrom() {
    document.getElementById("date_from").value = getDateFormat()
}

function getDateTo() {
    document.getElementById("date_to").value = getDateFormat()
}

// Get the information about boxes from the micro controller.
function makeSubmit() {

    target_from = document.getElementById("date_from");
    reqFrom = target_from.value + ", " + "00:01:00.000";

    target_to = document.getElementById("date_to");
    reqTo = target_to.value + ", " + "23:59:59.000";

    params = "?date_from=" + reqFrom

    params = params + "&date_to=" + reqTo


    alert(params)

    var ajaxRequest = new XMLHttpRequest();

    ajaxRequest.open("POST", "/getInfo", true);
    ajaxRequest.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    ajaxRequest.send(params);

}