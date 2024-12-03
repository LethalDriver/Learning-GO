from langchain.prompts import PromptTemplate
from langchain.prompts.few_shot import FewShotPromptTemplate
from templates import summary_example_template, summary_input_template
from multishot_examples import summary_examples, emotions_examples


def summarize_chat_multishot_query(chat_conversation: str) -> str:

    example_prompt = PromptTemplate(
        input_variables=["Conversation", "Summary"], template=summary_example_template
    )

    prompt = FewShotPromptTemplate(
        examples=summary_examples,
        example_prompt=example_prompt,
        suffix=summary_input_template,
        input_variables=["input"],
    )

    return prompt.format(input=chat_conversation)
