<template>
	<div class="instructions">
		<div class="d-flex justify-space-between my-2 py-4" @click="show = !show">
			<v-btn text class="text-caption" >Instructions</v-btn>
			<v-icon>{{ show ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
		</div>
		<v-expand-transition>
			<div v-show="show" class="py-2">
				<v-divider />
				<p>You can use this service with other tools, like curl!</p>

				<div class="notes" v-html="notes"></div>

				<h2>Fields</h2>
				<p>Only the primary field (marked) is required. the rest can be omitted entirely</p>

				<dl>
					<template v-for="f in fields">
						<dt>{{ f.name }}</dt>
						<dd><code>{{ f.type }}</code>: {{ f.description }}</dd>
					</template>

					<dt>burn</dt>
					<dd><code>boolean</code>: true if you want the message to self-destruct on read</dd>

					<dt>ttl</dt>
					<dd><code>int</code>: number of seconds for the message to exist. Afterwards, it self-deletes</dd>

					<dt>hidden</dt>
					<dd><code>boolean</code>: true if you want the message to be hidden from the listing. This means you need to keep track of the ID yourself</dd>
				</dl>


				<h2>Examples</h2>

				<div class="tool" v-for="t in tools" :key="t.name">
					<h2>{{ t.name }}</h2>

					<div class="action" v-for="a in t.actions" :key="`${t.name}-${a.name}`">
						<h3>{{ a.name }}</h3>

						<p v-if="a.note" v-html="a.note"></p>

						<div class="method" v-for="m,i in a.methods" :key="`${t.name}-${a.name}-${i}`">
							<h5 v-if="m.name" v-html="m.name"></h5>
							<pre v-html="m.code"></pre>
						</div>
					</div>
					<v-divider />
				</div>
			</div>
		</v-expand-transition>
	</div>
</template>


<script>
export default {
	props: {
		fields: {},
		tools: {},
		notes: {}
	},
	data() {
		return {
			show: false,
		}
	},

}
</script>

<style>

.instructions dt {
	font-weight: bold;
	font-weight: 800;
}
.instructions dd {
	padding-left: 20px;
}

.instructions dl {
	margin-bottom: 15px;
}

.tool {
	margin: 20px 0;
}

.action {
	padding: 20px 0;
}
.action h3 {
	margin-bottom: 10px;
}

.method {
	padding: 10px 0;
}

.method pre {
	padding-top: 5px;
	margin-left: 5px;
}

</style>