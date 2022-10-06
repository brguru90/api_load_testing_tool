const chart_option={
    chart: {
        id: "realtime",
        height: 350,
        type: "line",
        dropShadow: {
            enabled: true,
            color: "#000",
            top: 18,
            left: 7,
            blur: 10,
            opacity: 0.2,
        },
        animations: {
            enabled: true,
            easing: "easeout",
            speed: 1000,
            animateGradually: {
                enabled: true,
                delay: 150,
            },
            dynamicAnimation: {
                enabled: true,
                speed: 800,
            },
        },
        toolbar: {
            show: true,
            offsetX: -65,
            offsetY: 0,
            tools: {
                download: true,
                pan: true,
                selection: true,
                zoom: true,
                zoomin: true,
                zoomout: true,
                reset: true,
            },
            autoSelected: "zoom",
        },
        zoom: {
            enabled: true,
        },
    },
    dataLabels: {
        enabled: true,
    },
    stroke: {
        curve: "straight",
    },
    title: {
        text: 'API Response Time',
        align: "left",
        margin: 40,
        offsetX: 0,
        offsetY: -20,
        floating: false,
        style: {
            fontSize: "16px",
            fontWeight: "bold",
            fontFamily: undefined,
            color: "#263238",
        },
    },
    grid: {
        borderColor: "#e7e7e7",
        row: {
            colors: ["#f3f3f3", "transparent"], // takes an array which will be repeated on columns
            opacity: 0.5,
        },
    },
    markers: {
        size: 4,
        strokeWidth: 0,
    },
    xaxis: {
        type: "category",
        title: {
            text: "Iteration",
            style: {
                fontSize: "10px",
                color: "#78787D",
            },
        },
    },
    yaxis: {
        title: {
            text: "Time (ms)",
            style: {
                fontSize: "10px",
                color: "#78787D",
            },
        },
        labels: {
            rotate: -45,
        },
    },

    legend: {
        position: "top",
        horizontalAlign: "right",
        floating: true,
        offsetY: -25,
        offsetX: -5,
    },
}

export {
    chart_option
}