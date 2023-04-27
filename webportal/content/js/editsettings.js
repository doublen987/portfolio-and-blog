let myChart 
let myChart2
let reader = new FileReader();
let fileName = ""
let state = {
    timeChart: {
        year: 2023,
        month: 4
    },
    chronologicalData: [
            
    ]
}

function readURL(input) {
    return () => {
        if (input.files && input.files[0]) {
            reader = new FileReader();
            reader.onload = function (e) {
                let thumbnailImg = document.getElementById("settings-logo-image");
                thumbnailImg.src = e.target.result;
                thumbnailImg.width = 150;
                thumbnailImg.height = 200;
            };
            reader.readAsDataURL(input.files[0]);
        }
        fileName = input.files[0].name
    }
}

function onImageError() {
    //this.onerror=null;
    //if (this.src != '/content/no-image.png') {
        console.log("error!");
        console.log(this);
        this.src = '/content/no-image.png';
        this.width = 150;
        this.height = 200;
    //}
}

function saveSettings() {
    let websiteName = document.getElementById("settings-websitename").value
    let bgcolor1 = document.getElementById("settings-background-color-1").value
    let bgcolor2 = document.getElementById("settings-background-color-2").value
    let bgcolor3 = document.getElementById("settings-background-color-3").value
    let textcolor1 = document.getElementById("settings-text-color-1").value
    let textcolor2 = document.getElementById("settings-text-color-2").value
    let textcolor3 = document.getElementById("settings-text-color-3").value

    async function save() {
        const rawResponse = await fetch('/settings', {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            websiteName: websiteName,
            backgroundColor1: bgcolor1,
            backgroundColor2: bgcolor2,
            backgroundColor3: bgcolor3,
            textColor1: textcolor1,
            textColor2: textcolor2,
            textColor3: textcolor3,
            bytes: reader.result? reader.result.replace(new RegExp("data:image\/[^;]+;base64,"),'') : [],
            filename: fileName
        })
        });
        const content = await rawResponse;
      
        // console.log(content);
    };

    save()

    // console.log({
    //     websiteName: websiteName,
    //     backgroundColor1: bgcolor1,
    //     backgroundColor2: bgcolor2,
    //     backgroundColor3: bgcolor3,
    //     textColor1: textcolor1,
    //     textColor2: textcolor2,
    //     textColor3: textcolor3,
    //     bytes: reader.result? reader.result.replace(new RegExp("data:image\/[^;]+;base64,"),'') : []
    // })
}

function changeTimeChart(monthChange) {
    return () => {
        console.log(monthChange)
        let year = state.timeChart.year
        let maxyear = state.timeChart.year
        let month = state.timeChart.month + monthChange;
        let maxmonth = state.timeChart.month + monthChange + 1;
        // let newMin = year+ "-" + month + "-01"
        // let newMax = year+ "-" + (month + 1) + "-01"

        if(month > 12) {
            month = 1
            maxmonth = 2;
            maxyear++;
            year++;
        }
        if(month < 1) {
            month = 12
            maxmonth = 1;
            year--;
        }
        if(maxmonth > 12) {
            maxmonth = 1
            maxyear++;
        }
        if(maxmonth < 1) {
            maxmonth = 12
        }

        let newMin = year+ "-" + String(month).padStart(2,'0') + "-01"
        let newMax = maxyear+ "-" + String(maxmonth).padStart(2,'0') + "-01"
        
        console.log(newMin)
        console.log(newMax)

        myChart.options.scales.x.min = newMin
        myChart.options.scales.x.max = newMax
        myChart.update();
        state = {
            ...state,
            timeChart: {
                ...state.timeChart,
                month: month,
                year: year
            }
        }
        console.log(state)
    }
}

function initCharts() {
    const ctx = document.getElementById('myChart');
    const ctx2 = document.getElementById('myChart2');
    myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            datasets: [{
                label: '# of site visitors per day',
                data: [
                    // {
                    //     x: '2022/04/06',
                    //     y: 50
                    // }, 
                    // {
                    //     x: '2022/08/07',
                    //     y: 60
                    // }, 
                    {
                        x: '2022-11-10',
                        y: 20
                    }
                ],
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                borderColor: 'rgba(255, 99, 132, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'day',
                        displayFormats: {
                            day: 'yyyy-MM-dd'

                        }
                    },
                    min: '2023-04-01',
                    max: '2023-04-30'
                },
                y: {
                    beginAtZero: true
                }
            }
        }
    });

    const data2 = {
        labels: [
            'Red',
            'Blue',
            'Yellow'
        ],
        datasets: [{
            label: 'My First Dataset',
            data: [300, 50, 100],
            backgroundColor: [
            'rgb(255, 99, 132)',
            'rgb(54, 162, 235)',
            'rgb(255, 205, 86)'
            ],
            hoverOffset: 4
        }]
    };

    const config2 = {
        type: 'doughnut',
        data: data2,
    };


    myChart2 = new Chart(ctx2, config2);
    console.log("chart initialized")
}

async function getStats() {
    let response = await fetch('/stats', {
        method: 'GET',
        credentials: 'include'
    })
    let data = await response.json()
    console.log(data)
    let newData =[]
    for(let i = 0; i < data.length; i++) {
        let pageViews = 0
        let pageNames = Object.keys(data[i].pages)
        for(let j = 0; j < pageNames.length; j++) {
            pageViews += data[i].pages[pageNames[j]]
        }
        let dateArray = data[i].date.split("/")
        newData.push({
            x: dateArray[0] + "-" + dateArray[1] + "-" + dateArray[2],
            y: pageViews
        })
    }
    console.log(data)

    let countries = {}

    for(let i = 0; i < data.length; i++) {
        let pageNames = Object.keys(data[i].countries)
        for(let j = 0; j < pageNames.length; j++) {
            if(pageNames[j] != "InvalidIPaddress" && pageNames[j] != "-") {
                if(countries[pageNames[j]] != undefined)
                    countries[pageNames[j]] += data[i].countries[pageNames[j]]
                else
                    countries[pageNames[j]] = 0
            }
        }
    }

    myChart.data.datasets[0].data = newData

    console.log(countries)

    myChart2.data.labels = Object.keys(countries)
    myChart2.data.datasets[0].data = Object.keys(countries).map(countryName => {
        return countries[countryName] 
    })
    
    myChart.update();
    myChart2.update();
    console.log(myChart.data)
}

window.onload = () => {

    var thumbnailImg = document.getElementById("settings-logo-image");
    console.log("bla")
    thumbnailImg.width = 150;
    thumbnailImg.height = 200;
    thumbnailImg.onerror = onImageError;

    let thumbnailInput = document.getElementById("settings-logo-input");
    thumbnailInput.onchange = readURL(thumbnailInput);


    let saveSettingsButton = document.getElementById("save-settings-button");
    saveSettingsButton.addEventListener("click", () => {
        saveSettings()
    })

    initCharts()

    let stats = getStats()

    let buttonNext = document.getElementById("button-next");
    let buttonPrev = document.getElementById("button-prev");
    buttonNext.addEventListener("click", changeTimeChart(+1))
    buttonPrev.addEventListener("click", changeTimeChart(-1))
}