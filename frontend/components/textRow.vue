<template>
	<v-row>
		<v-col cols="2">
			<v-icon v-if="burn">mdi-fire</v-icon>
			<v-chip v-if="ttlExp">
				<v-icon left :title="ttl">mdi-alarm</v-icon>
				{{ timeLeft }}
			</v-chip>
		</v-col>
		<v-col>
			<nuxt-link :to="'/text/'+id">{{ !burn ? text : '&lt;censored&gt;' }}</nuxt-link>
		</v-col>
	</v-row>	
</template>


<script>
import { formatDistanceToNow } from 'date-fns'

export default {
	props: {
		id: {},
		ttl: {},
		burn: {},
		text: {},
		created: {},
	},
	data() {
		return {

		}
	},
	computed: {
		ttlExp() {
			if (this.created && this.ttl) {
				return new Date( (this.created+this.ttl) * 1000 )
			}
			return null
		},
		timeLeft() {
			if (this.ttlExp) {
				return formatDistanceToNow(
					this.ttlExp,
					{ addSuffix: true, includeSeconds: true }
				)
			}
			return '';
		}
	}
}
</script>