<script lang="ts">
  import QuizCard from "../../lib/QuizCard.svelte";
  import type { Quiz } from "../../model/quiz";


    let quizzes:Quiz[]=[]

    async function getQuizes() {
    try {
      const response = await fetch("http://localhost:5001/api/quizzes");

      let json = await response.json();
      console.log(json);
      quizzes = json;
    } catch (error) {
      console.log("Error",error);
    }
  }

  getQuizes()




</script>
<div class="p-8">
    <h2 class="text-4xl font-bold">Your quizzes</h2>
    <div class="flex flex-col gap-2 mt-4">
    {#each quizzes as quiz(quiz.id)}
        <QuizCard on:host {quiz} />
    {/each}
    </div>
  </div>

