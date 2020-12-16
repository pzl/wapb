<template>
	<v-row dense :class="{ burned: burn, expirable: ttlExp }">
		<v-col cols="auto">
			<v-chip small v-if="ttlExp" :title="expireDateText">
				<v-icon left>mdi-alarm</v-icon>
				{{ timeLeft }}
			</v-chip>
			<v-icon color="orange" v-if="burn">mdi-fire</v-icon>
		</v-col>
		<v-col cols="auto"><nuxt-link :to="'/file/'+id">{{ id }}</nuxt-link></v-col>
		<v-col class="previewText">{{ !burn ? fileInfo : '&lt;censored&gt;' }}</v-col>
		<v-col class="fileDetails">{{ fileDetails }}</v-col>
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
import fileSize from '~/helpers/fileSize'

export default {
	props: {
		id: {},
		ttl: {},
		burn: {},
		files: {},
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
		fileInfo() {
			if (!this.files) {
				return ''
			}
			return `${this.files.length} File${this.files.length > 1 ? 's' : ''}`
		},
		fileDetails() {
			if (this.burn || !this.files) {
				return ''
			}
			if (this.files.length > 2) {
				return "Total: "+fileSize(this.files.reduce((acc, f) => acc + f.size, 0))
			}
			return this.files.map(f => `${f.name}, ${fileSize(f.size)}`).join('; ')
		}
	}
}
</script>

<style>

.fileDetails {
	text-align: right;
}

.burned .previewText {
	font-family: monospace;
}

</style>