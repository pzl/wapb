<template>
	<v-row>
		<v-col cols="12" md="">
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-textarea name="text" label="Text" hint="New Text" v-model="toCreate.text" outlined />
			<create-meta
				ref="meta"
				v-bind.sync="toCreate"
			/>
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="toCreate.text == ''">Create</v-btn>
			<instructions v-bind="instructions" />
		</v-col>
		<v-col cols="12" md="">
			<text-row v-for="t in texts" :key="t.id" v-bind="t" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import TextRow from '~/components/textRow'
import Instructions from '~/components/instructions'

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
			texts: [],


			instructions: {
				fields: [{ name: "text", type: "string", description: "(required) Content of the paste" }],
				tools: [
					{
						name: "Curl",
						actions: [
							{
								name: "Create",
								methods: [
									{ name: "Simple", code: `curl -d "Some text here" ${ this.$server }/api/v1/text` },
									{ name: "Simple with URL parameters", code: `curl -d "Some text here" "${ this.$server }/api/v1/text?burn=true&ttl=3600"` },
									{ name: "Using Form Parameters", code: `curl -d "text=text body content here" -d "burn=true" -d "ttl=3600" ${ this.$server }/api/v1/text` },
									{ name: "File <em>contents</em> to text snippet", code: `curl --data-binary @filename ${ this.$server }/api/v1/text` },
									{ name: "From your clipboard", code: `pbpaste | curl --data-binary @- ${ this.$server }/api/v1/text` },
									{ name: "JSON", code: `curl -d '{ "text": "my message", "ttl":86500, "burn": false }' ${ this.$server }/api/v1/text` },
								]
							},
							{
								name: "Retrieve",
								note: `First, create an entry as above. When created successfully, you will get an ID back. It will either be at the end of a URL: (<code>${ this.$server }/text/fe4c71</code>) or as a JSON message: <code>{"id":"02b268", ... }</code>. We will use ID XXYY2 in our examples`,
								methods: [
									{ name: "As plain text", code: `curl -H "Accept: text/plain" ${ this.$server }/api/v1/text/XXYY2` },
									{ name: "To your clipboard!", code: `curl -H "Accept: text/plain" ${ this.$server }/api/v1/text/XXYY2 | pbcopy` },
									{ name: "As JSON", code: `curl ${ this.$server }/api/v1/text/XXYY2` },
								]
							},
							{ name: "Delete", methods: [{ code: `curl -X DELETE ${ this.$server }/api/v1/text/XXYY2` }] },
						]
					},
					{
						name: "HTTPie",
						actions: [
							{
								name: "Create",
								methods: [
									{ code: `http POST ${ this.$server }/api/v1/text 'text=whatever you want'` },
									{ code: `http POST ${ this.$server }/api/v1/text 'text=whatever you want' burn:=true ttl:=3600` },
									{ code: `echo '{ "text":"explicit JSON object", "burn":true }' | http POST ${ this.$server }/api/v1/text` },
								]
							},
							{ name: "Retrieve", methods: [{ code: `http ${ this.$server }/api/v1/text/XXYY2` }] },
							{ name: "Delete", methods: [{ code: `http DELETE ${ this.$server }/api/v1/text/XXYY2` }] },
						]
					}
				]
			}
		}
	},
	async asyncData({ $http, $server }) {
		const texts = await $http.$get(`${$server}/api/v1/text`).then(d => d.data)
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
			const data = await this.$http.$post(`${this.$server}/api/v1/text`, this.toCreate).then(d => {
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
	components: { CreateMeta, TextRow, Instructions }

}
</script>