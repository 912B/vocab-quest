export const Audio = {
    speak(text) {
        if ('speechSynthesis' in window) {
            const utterance = new SpeechSynthesisUtterance(text);
            // Try to set voice?
            utterance.rate = 1.0;
            window.speechSynthesis.speak(utterance);
        }
    },

    play(soundName) {
        console.log("Play Sound:", soundName);
    }
};
