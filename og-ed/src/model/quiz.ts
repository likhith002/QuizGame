 interface Quiz {
  id: string;
  name: string;
  questions: QuizQuestion;
}

 interface QuizQuestion {
  id: string;
  name: string;
  choices: QuizChoice[];
}

 interface QuizChoice {
  id: string;
  name: string;
  correct: boolean;
}

 interface Player {
  id: string;
  name: string;
}
