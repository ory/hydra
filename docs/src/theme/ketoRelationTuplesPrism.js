const delimiter = {
  delimiter: /[:#@()]/
}

const namespace = {
  pattern: /[^:#@()\n]+:/,
  inside: {
    ...delimiter,
    namespace: /.*/
  }
}

const object = {
  pattern: /[^:#@()\n]+#/,
  inside: {
    ...delimiter,
    'property-access': /.*/
  }
}

const relation = {
  pattern: /[^:#@()\n]+/
}

const subjectID = {
  pattern: /@[^:#@()\n]+/,
  inside: {
    ...delimiter,
    subject: /.*/
  }
}

const subjectSet = {
  pattern: /@\(([^:#@()\n]+:)?([^:#@()\n]+)#([^:#@()\n]*)\)/,
  inside: {
    delimiter: /[@()]*/,
    namespace,
    object,
    relation
  }
}

export default (prism) =>
  (prism.languages['keto-relation-tuples'] = {
    comment: /\/\/.*(\n|$)/,
    'relation-tuple': {
      pattern:
        /([^:#@()\n]+:)?([^:#@()\n]+)#([^:#@()\n]+)@?((\(([^:#@()\n]+:)?([^:#@()\n]+)#([^:#@()\n]*)\))|([^:#@()\n]+))/,
      inside: {
        namespace,
        object,
        subjectID,
        subjectSet,
        relation
      }
    }
  })
