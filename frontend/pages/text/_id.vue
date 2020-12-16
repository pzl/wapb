<template>
	<v-container>
		<v-alert v-if="alert.show" :type="alert.type">{{ alert.text }}</v-alert>

		<v-alert icon="mdi-fire" type="warning" v-if="burn">This message has self-destructed. Make sure you know its contents. It will be gone if you try to refresh or close the tab</v-alert>

		<v-card shaped :loading="loading">
			<template slot="progress">
				<v-progress-linear color="deep-purple" height="10" indeterminate />
			</template>

			<v-card-title>
				<v-row>
					<v-col>{{ id }}</v-col>
					<v-col v-if="burn" class="text-right"><v-icon large color="red">mdi-fire</v-icon></v-col>
				</v-row>
			</v-card-title>
			<v-card-subtitle>{{ creationTime }}</v-card-subtitle>
			<v-card-text class="content-text">
				<v-textarea outlined hide-details full-width v-if="showAsEditable" :value="text" />
				<pre v-else>{{ text }}</pre>
			</v-card-text>
			<v-card-actions>
				<v-btn icon @click="showAsEditable = !showAsEditable"><v-icon>mdi-file-document-edit-outline</v-icon></v-btn>
				<v-spacer />
				<v-btn color="red" text @click="processDelete">Delete</v-btn>
			</v-card-actions>
		</v-card>

	</v-container>
</template>


<script>
import { format } from 'date-fns';

export default {
	data () {
		return {
			id: "",
			text: "",
			burn: false,
			ttl: null,
			created: 0,


			loading: false,
			showAsEditable: false,
			alert: {
				text: '',
				show: false,
				type: '',
			}
		}
	},
	async asyncData(context) {
		const data =  await context.$http.$get(`${context.$server}/api/v1/text/${context.params.id}`)
							.catch(e => {
								console.log(e)
								context.error(e)
							})
		return data
	},
	methods: {
		async processDelete() {
			this.loading = true;
			this.alert.show = false;
			await this.$http.$delete(`${this.$server}/api/v1/text/${this.id}`)
							.then(() => {
								this.$router.push({
									path: "/text"
								})
							})
							.catch(e => {
								this.loading = false;
								console.log(e)
								this.alert.text = e.message;
								this.alert.type = "error"
								this.alert.show = true;
							})
		}
	},
	computed: {
		creationTime() {
			return format(new Date(this.created*1000), 'EE PPpp ')
		}
	}
}
</script>


<style>

.content-text pre {
	white-space: pre-wrap;
}
</style>