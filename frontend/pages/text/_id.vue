<template>
	<div>
		<h1>{{ text }}</h1>

		<p v-if="burn">This message has self-destructed. Make sure you know its contents. It will be gone if you try to refresh or close the tab</p>
	</div>
</template>


<script>
export default {
	data () {
		return {
			text: "",
			burn: false,
		}
	},
	async asyncData(context) {
		const data =  await context.$http.$get('http://localhost:7473/api/v1/text/'+context.params.id)
							.catch(e => {
								console.log(e)
								context.error(e)
							})
		return data
	}
}
</script>