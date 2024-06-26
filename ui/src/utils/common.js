import { useToast } from '@/components/ui/toast/use-toast'

const { toast } = useToast()

const categoryColors = {
  food: 'bg-green-500',
  entertainment: 'bg-blue-500',
  groceries: 'bg-yellow-500',
  clothing: 'bg-orange-500',
  travel: 'bg-purple-500',
  utilities: 'bg-red-500',
  shopping: 'bg-pink-500',
  health: 'bg-indigo-500',
  education: 'bg-teal-500'
}

export function getCategoryColor(category) {
  const lowerCaseCategory = category.toLowerCase()
  return categoryColors[lowerCaseCategory] || 'bg-gray-500' // default color if category is not found
}

export const showToast = (title, description, isError = false) => {
  toast({
    title: title,
    description: description,
    variant: isError ? 'destructive' : 'default'  // Use 'destructive' for errors, 'default' otherwise
  });
};