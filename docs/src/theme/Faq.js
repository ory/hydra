import React, { useState } from 'react'
import cn from 'classnames'
import styles from './faq.module.css'
import genStyle from './faq.gen.module.css'

const Question = ({ children, tags }) => (
  <div className={cn(genStyle.question, ...tags.map((tag) => genStyle[tag]))}>
    {children}
  </div>
)

const TagButton = ({ tag, isSelected, children, toggleSelected }) => (
  <li
    className={cn(
      { [genStyle.selected]: isSelected },
      genStyle[tag],
      styles.pills,
      styles.pills__item,
      { [styles['pills__item--active']]: isSelected }
    )}
    onClick={toggleSelected}
  >
    {children}
  </li>
)

const FaqTags = ({ tags, initiallyDisabled }) => {
  const [selectedTags, setSelectedTags] = useState(
    tags.filter((t) => !initiallyDisabled.includes(t))
  )

  return (
    <>
      {tags.map((tag) => (
        <TagButton
          key={tag}
          tag={tag}
          isSelected={selectedTags.find((t) => t === tag)}
          toggleSelected={() => {
            if (selectedTags.find((t) => t === tag)) {
              setSelectedTags(selectedTags.filter((t) => t !== tag))
            } else {
              setSelectedTags([...selectedTags, tag])
            }
          }}
        >
          #{tag}
        </TagButton>
      ))}
    </>
  )
}

export { FaqTags, Question }
