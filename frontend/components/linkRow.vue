<template>
	<v-row dense :class="{ burned: burn, expirable: ttlExp }">
		<v-col cols="auto">
			<v-chip small v-if="ttlExp" :title="expireDateText">
				<v-icon left>mdi-alarm</v-icon>
				{{ timeLeft }}
			</v-chip>
			<v-icon color="orange" v-if="burn">mdi-fire</v-icon>
		</v-col>
		<v-col cols="auto"><nuxt-link :to="'/link/'+id">{{ id }}</nuxt-link></v-col>
		<v-col class="previewText">
			<template v-if="burn">&lt;censored&gt;</template>
			<a v-else :href="url">{{ truncateURL  }}</a>
		</v-col>
		<!--
		<v-col justify="end" class="text-right">
			<v-btn icon small @click="">
				<v-icon dense color="grey lighten-2">mdi-delete</v-icon>
			</v-btn>
		</v-col>
		-->
	</v-row>	
</template>


<script>
import { formatDistanceToNowStrict, format } from 'date-fns'

export default {
	props: {
		id: {},
		ttl: {},
		burn: {},
		url: {},
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
				return formatDistanceToNowStrict(
					this.ttlExp,
					{ includeSeconds: true }
				)
			}
			return '';
		},
		expireDateText() {
			return format(this.ttlExp, 'EEE PPpp')
		},
		truncateURL() {
			if (this.url.length > 25) {
				let trunc = this.url.slice(0,40)
				let lastSpace = trunc.lastIndexOf(" ")
				if (lastSpace == -1) {
					lastSpace = 25
				}
				trunc = trunc.slice(0,lastSpace)
				return trunc + "…"
			}
			return this.url;
		}
	}
}
</script>

<style>

.burned .previewText {
	font-family: monospace;
}

</style>