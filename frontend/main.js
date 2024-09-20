const BASE_URL = 'http://localhost:8080';
const audioTrack = document.getElementById('audioTrack');
const oracleText = document.getElementById('oracleText');
const visualizer = document.getElementById('visualizer');
const startButton = document.getElementById('startButton');

const audioContext = new (window.AudioContext || window.webkitAudioContext)();
const analyser = audioContext.createAnalyser();
analyser.fftSize = 256;
const dataArray = new Uint8Array(analyser.frequencyBinCount);

let source;

function createVisualizer() {
    for (let i = 0; i < 50; i++) {
        visualizer.appendChild(document.createElement('div')).className = 'visualizer-circle';
    }
}

function updateVisualizer() {
    analyser.getByteFrequencyData(dataArray);
    const circles = visualizer.children;
    const centerX = visualizer.clientWidth / 2;
    const centerY = visualizer.clientHeight / 2;
    const maxRadius = Math.min(centerX, centerY);

    for (let i = 0; i < circles.length; i++) {
        const value = dataArray[i] / 255;
        const size = value * maxRadius * 2;
        const opacity = value * 0.5 + 0.1;
        circles[i].style.cssText = `
            width: ${size}px;
            height: ${size}px;
            left: ${centerX - size / 2}px;
            top: ${centerY - size / 2}px;
            opacity: ${opacity};
        `;
    }
    requestAnimationFrame(updateVisualizer);
}

function playAudio() {
    oracleText.textContent = "Processing Query...";
    audioTrack.play().then(() => {
        if (!source) {
            source = audioContext.createMediaElementSource(audioTrack);
            source.connect(analyser);
            analyser.connect(audioContext.destination);
        }
        updateVisualizer();
    }).catch(error => console.error('Error initializing audio:', error));

    audioTrack.onended = () => oracleText.textContent = "Ready for Next Query";
}

function checkForNewMedia() {
    fetch(`${BASE_URL}/getLatestAudio`)
        .then(response => response.json())
        .then(data => {
            audioTrack.src = data.url;
            audioTrack.load();
            playAudio();
        })
        .catch(error => console.error('Error fetching new audio URL:', error));
}

startButton.addEventListener('click', () => {
    audioContext.resume().then(() => {
        createVisualizer();
        checkForNewMedia();
        setInterval(checkForNewMedia, 10000);
    });
});