//@ts-check

///<reference path="../htmx/htmx.d.ts" />

/**
 * @typedef Word
 * @property {string} id
 * @property {string} word
 * @property {string} translations
 * @property {number} stageLevel
 * @property {number | undefined} wrongAnswers
 * @property {string} stageName
 */

const STAGE_COLORS = {
  apprentice: "#db2777",
  guru: "#9333ea",
  master: "#2563eb",
  enlightened: "#0284c7",
  burned: "#57534e"
};

function main() {
  const { currentWordEl, checkBtn, guessInput, reviewTrigger, wordSection } =
    getElements();

  let words = getWordsData();
  let guess = "";

  const currentWord = createSignal(words[0]);
  currentWord.subscribe((v) => {
    currentWordEl.innerText = v.word;
    wordSection.style.backgroundColor = STAGE_COLORS[v.stageName];
  });

  const correctGuessed = createSignal(/** @type {null | boolean} */ (null));
  correctGuessed.subscribe((correct) => {
    guessInput.disabled = correct !== null;

    if (correct === null) {
      guessInput.classList.remove("bg-red-500", "bg-green-500");
      return;
    }

    guessInput.classList.toggle("bg-red-500", !correct);
    guessInput.classList.toggle("bg-green-500", correct);
  });

  guessInput.addEventListener("change", (e) => {
    const target = /** @type HTMLInputElement | null */ (e.currentTarget);
    if (!target) return;

    guess = target.value;
  });

  checkBtn.addEventListener("click", () => {
    if (correctGuessed.value === null) {
      correctGuessed.value = checkWord(currentWord.value, guess);
      return;
    }

    const correct = correctGuessed.value;
    const wrongAnswers = currentWord.value.wrongAnswers ?? 0;

    correctGuessed.value = null;
    guessInput.value = "";

    if (!correct) {
      currentWord.value.wrongAnswers = wrongAnswers + 1;
      words = [...words.slice(1), currentWord.value];
      currentWord.value = words[0];

      return;
    }

    words = words.slice(1);
    const shouldFinish = words.length === 0;
    reviewTrigger.setAttribute(
      "hx-vals",
      JSON.stringify({
        id: currentWord.value.id,
        wrongAnswers,
        finish: shouldFinish
      })
    );

    htmx.trigger(reviewTrigger, "word-review");

    if (!shouldFinish) {
      currentWord.value = words[0];
    }
  });
}

function getWordsData() {
  const dataElement = document.getElementById("review-words-data");
  if (!dataElement || !dataElement.textContent) {
    throw "no words data found";
  }

  /** @type Word[] */
  const words = JSON.parse(dataElement.textContent);
  dataElement.remove();
  return words;
}

function getElements() {
  const currentWordEl = getElement("current-word");
  const checkBtn = getElement("check-button");
  const guessInput = /** @type HTMLInputElement */ (getElement("guess-input"));
  const reviewTrigger = getElement("review-trigger");
  const wordSection = getElement("word-section");

  return { currentWordEl, checkBtn, guessInput, reviewTrigger, wordSection };
}

/**
 *
 * @param {string} id
 */
function getElement(id) {
  const elem = document.getElementById(id);
  if (!elem) {
    throw `no element with id "${id}" found`;
  }

  return elem;
}

/**
 *
 * @param {Word} word
 * @param {string} guess
 * @returns {boolean}
 */
function checkWord(word, guess) {
  const translations = word.translations
    .split(",")
    .map((t) => t.trim().toLowerCase());
  return translations.includes(guess.trim().toLowerCase());
}

/**
 * @template {unknown} T
 * @param {T} initialValue
 */
function createSignal(initialValue) {
  let _value = initialValue;

  /** @typedef {((value: T) => void)} Subscriber  */

  /** @type Subscriber[] */
  let subscribers = [];

  function notify() {
    for (let subscriber of subscribers) {
      subscriber(_value);
    }
  }

  return {
    get value() {
      return _value;
    },
    set value(v) {
      _value = v;
      notify();
    },

    /**
     * @param {Subscriber} subscriber
     */
    subscribe: (subscriber) => {
      subscriber(_value);
      subscribers.push(subscriber);
    }
  };
}

main();
