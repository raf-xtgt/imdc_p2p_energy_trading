function submitInput() {    
    let inputText = document.getElementById("textsubmit").value;
    let myURL = "http://localhost:8080/WriteBlock"
    var xhr = new XMLHttpRequest();
    xhr.open("POST", myURL, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.add
    xhr.send(JSON.stringify({
        "BPM": inputText
    }));

    /* let headers = new Headers();

    headers.append('Content-Type','application/json');
    headers.append('Accept','application/json');

    headers.append('Access-Control-Allow-Origin','*');
    headers.append('GET','POST','OPTIONS');

    console.log(inputText) */
}
