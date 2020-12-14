<template>
	<div>
		<v-card shaped :loading="loading">
			<template slot="progress">
				<v-progress-linear color="deep-purple" height="10" indeterminate />
			</template>

			<v-card-title>{{ id }}</v-card-title>
			<v-card-subtitle>{{ creationTime }}</v-card-subtitle>
			<v-card-text>
				<pre>{{ text }}</pre>
			</v-card-text>
			<v-card-actions>
				<v-btn color="red" text @click="processDelete">Delete</v-btn>
			</v-card-actions>
		</v-card>

		<v-alert type="warning" v-if="burn">This message has self-destructed. Make sure you know its contents. It will be gone if you try to refresh or close the tab</v-alert>
	</div>
</template>


<script>
export default {
	data () {
		return {
			id: "",
			text: "",
			burn: false,
			ttl: null,
			created: 0,
			loading: false,
		}
	},
	async asyncData(context) {
		const data =  await context.$http.$get('http://localhost:7473/api/v1/text/'+context.params.id)
							.catch(e => {
								console.log(e)
								context.error(e)
							})
		return data
	},
	methods: {
		processDelete() {
			this.loading = true;
			setTimeout(() => { this.loading = false; }, 1500)
		}
	},
	computed: {
		creationTime() {
			return new Date(this.created*1000)
		}
	}
}
</script>