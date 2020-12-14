<template>
	<v-row>
		<v-col>
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-textarea name="text" label="Text" hint="New Text" v-model="toCreate.text" outlined />
			<create-meta
				ref="meta"
				v-bind.sync="toCreate"
			/>
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="toCreate.text == ''">Create</v-btn>
		</v-col>
		<v-col>
			<text-row v-for="t in texts" :key="t.id" v-bind="t" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import TextRow from '~/components/textRow'

const createFactory = () => {
	return {
		text: '',
		burn: false,
		ttl: null,
		hidden: false,
	}
} 

export default {
	data () {
		return {
			toCreate: createFactory(),
			isLoading: false,
			alert: null,
			texts: []
		}
	},
	async asyncData() {
		const texts = await fetch('http://localhost:7473/api/v1/text').then(res => res.json()).then(d => d.data)
		return { texts }
	},
	methods: {
		async create() {
			if (this.toCreate.text == '') {
				this.err = { message: 'Cannot create empty text' };
				return
			}
			this.alert = null;
			this.isLoading = true;
			const data = await this.$http.$post('http://localhost:7473/api/v1/text', this.toCreate).then(d => {
				this.isLoading = false;
				this.toCreate = createFactory()
				this.$refs.meta.hideTTL()
				this.texts.push(d)
				this.alert = {
					type: "success",
					message: `Created ${d.id}`,
				}
			}).catch(e => {
				this.isLoading = false;
				this.alert = {
					type: "error",
					message: e.message,
				}
			})
		},
	},
	components: { CreateMeta, TextRow }

}
</script>