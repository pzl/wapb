<template>
	<div>
		<h1>Links</h1>

		<div v-for="l in links" :key="l.id">
			<p><a :href="l.dest">{{ l.dest }}</a></p>
		</div>
	</div>
</template>


<script>
export default {
	data () {
		return {
			links: []
		}
	},
	async asyncData() {
		const links = await fetch('https://localhost:7473/api/v1/links').then(res => res.json()).then(d => d.data)
		return { links }
	}
}
</script>