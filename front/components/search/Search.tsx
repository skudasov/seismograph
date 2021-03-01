import React, { HTMLAttributes, useState } from 'react'
import TextField from '@material-ui/core/TextField'
import Autocomplete from '@material-ui/lab/Autocomplete'
import { ConnectedProps, connect } from 'react-redux'
import { Dispatch } from 'redux'
// eslint-disable-next-line @typescript-eslint/ban-ts-ignore
// @ts-ignore
// import { v4 as uuidv4 } from 'uuid'
import { RootState } from '../../reducers'

const top100Films = [
  { title: 'The Shawshank Redemption', year: 1994 },
  { title: 'The Godfather', year: 1972 },
  { title: 'The Godfather: Part II', year: 1974 },
  { title: 'The Dark Knight', year: 2008 },
  { title: '12 Angry Men', year: 1957 },
  { title: "Schindler's List", year: 1993 },
]

const mapState = (state: RootState) => ({
  app: state.app,
})

const mapDispatch = (dispatch: Dispatch) => ({})

const connector = connect(mapState, mapDispatch)
type SearchBarProps = ConnectedProps<typeof connector> &
  HTMLAttributes<HTMLDivElement>

const SearchBar = (props: SearchBarProps) => {
  const [inputData, setInputData] = useState('')

  // eslint-disable-next-line no-shadow
  const { style } = props
  return (
    <Autocomplete
      id="combo-box-demo"
      options={top100Films}
      getOptionLabel={(option) => option.title}
      style={style}
      renderInput={(params) => (
        <TextField
          {...params}
          data-test="add-word-input"
          onKeyDown={(e: React.KeyboardEvent) => {
            if (e.key === 'Enter') {
              // setWord(inputData)
              // setWordCard({ title: uuidv4() })
            }
          }}
          onChange={(e) => {
            setInputData(e.target.value)
          }}
          label="Search test"
          variant="filled"
        />
      )}
    />
  )
}

export default connector(SearchBar)
