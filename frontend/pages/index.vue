<template>
  <div>
    <h1>Create Something</h1>
  </div>
</template>

<script>

export default {
  middleware({ redirect }) {
    return redirect('/text')
  }
}
</script>
