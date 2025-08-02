import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/core/ui/label_item_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/card_details_widgets/priority_selector.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/card_details_widgets/status_widget.dart';
import 'package:swiftly_mobile/ui/kanban/widgets/date_widget.dart';

import '../../../../domain/kanban/models/card_item.dart';
import '../../../../domain/kanban/models/priority.dart';
import '../../../../domain/models/label_item.dart';
import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/label_notifier_provider.dart';
import '../../../../providers/user_notifier_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import '../../../home/widgets/content/avatar_widget.dart';
import 'button_widget.dart';
import 'text_field_widget.dart';

class CartDetailsWidget extends ConsumerStatefulWidget {
  final CardItem card;

  const CartDetailsWidget({super.key, required this.card});

  @override
  ConsumerState<CartDetailsWidget> createState() => _CartDetailsWidgetState();
}

class _CartDetailsWidgetState extends ConsumerState<CartDetailsWidget> {
  late final TextEditingController _titleController;
  late final TextEditingController _descriptionController;
  late String _selectedColumnId;
  Priority? _selectedPriority;

  @override
  void initState() {
    super.initState();
    _titleController = TextEditingController(text: widget.card.title);
    _descriptionController = TextEditingController(
      text: widget.card.description,
    );
    _selectedColumnId = widget.card.columnId;
    _selectedPriority = widget.card.priority;
  }

  @override
  Widget build(BuildContext context) {
    final users = ref.watch(userNotifierProvider).users;
    final currentUser = users.firstWhere(
      (user) => user.id == widget.card.userId,
    );
    final labels =
        ref
            .watch(labelNotifierProvider)
            .labels
            .where((label) => label.cardId == widget.card.id)
            .toList();

    return Center(
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
          child: SingleChildScrollView(
            child: Container(
              constraints: const BoxConstraints(maxWidth: 400),
              width: 400,
              height: 500,
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
              decoration: BoxDecoration(
                color: AppColors.white26,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      if (labels.isNotEmpty)
                        Row(
                          children: [
                            ...labels.map(
                              (l) => Stack(
                                children: [
                                  Padding(
                                    padding: const EdgeInsets.only(right: 8),
                                    child: LabelItemWidget(labelItem: l),
                                  ),
                                  Positioned(
                                    top: 0,
                                    right: 0,
                                    child: IconButton(
                                      onPressed: () => _handleRemoveLabel(l.id),
                                      icon: const Icon(Icons.clear, size: 16),
                                      color: Colors.blue,
                                      padding: EdgeInsets.zero,
                                      constraints: const BoxConstraints(),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            if (labels.length < 5)
                              AddLabelButton(
                                onTap: () => _handleAddLabel(widget.card.id),
                              ),
                          ],
                        )
                      else
                        AddLabelButton(
                          onTap: () => _handleAddLabel(widget.card.id),
                        ),

                      const Spacer(),
                      DateWidget(
                        date: widget.card.createdAt,
                        color: AppColors.twitter,
                      ),
                    ],
                  ),
                  const SizedBox(height: 10),
                  SizedBox(
                    width: 300,
                    child: TextFieldWidget(
                      hintText: 'Название карточки',
                      controller: _titleController,
                    ),
                  ),
                  const SizedBox(height: 10),
                  SizedBox(
                    width: 300,
                    child: TextFieldWidget(
                      hintText: 'Описание',
                      controller: _descriptionController,
                    ),
                  ),
                  const SizedBox(height: 10),
                  Row(
                    children: [
                      const Text(
                        'Статус',
                        style: TextStyle(color: AppColors.white),
                      ),
                      const SizedBox(width: 10),
                      DropDownMenuWidget(
                        statusId: _selectedColumnId,
                        onChanged: (newId) {
                          setState(() {
                            _selectedColumnId = newId;
                          });
                        },
                      ),
                    ],
                  ),
                  const SizedBox(height: 10),
                  if (_selectedPriority != null)
                    PrioritySelector(
                      selected: _selectedPriority!,
                      onSelected: (priority) {
                        setState(() {
                          _selectedPriority = priority;
                        });
                      },
                    ),
                  const SizedBox(height: 20),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      const Text(
                        'Создатель задачи',
                        style: TextStyle(color: AppColors.white),
                      ),
                      const SizedBox(width: 10),
                      AvatarWidget(imageUrl: currentUser.image),
                      const SizedBox(width: 10),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              '${currentUser.name} ${currentUser.lastName ?? ''}',
                              style: AppTextStyles.style13,
                              overflow: TextOverflow.ellipsis,
                              maxLines: 1,
                              softWrap: false,
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const Spacer(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      ButtonWidget(
                        title: 'Удалить',
                        color: AppColors.red,
                        onTap: () => _handleDelete(context),
                      ),
                      const SizedBox(width: 10),
                      ButtonWidget(
                        title: 'Сохранить',
                        color: AppColors.blue,
                        onTap: () => _handleSave(context),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _handleDelete(BuildContext context) {
    ref.read(cardNotifierProvider.notifier).removeCard(widget.card.id);
    Navigator.of(context).pop();
  }

  void _handleSave(BuildContext context) {
    final newTitle = _titleController.text;
    final newDescription = _descriptionController.text;
    ref
        .read(cardNotifierProvider.notifier)
        .updateTitle(widget.card.id, newTitle);
    ref
        .read(cardNotifierProvider.notifier)
        .updateDescription(widget.card.id, newDescription);
    ref
        .read(cardNotifierProvider.notifier)
        .updateColumn(widget.card.id, _selectedColumnId);
    ref
        .read(cardNotifierProvider.notifier)
        .updatePriority(widget.card.id, _selectedPriority!);
    Navigator.of(context).pop();
  }

  void _handleAddLabel(String cardId) {
    final newLabel = LabelItem.create(
      cardId: cardId,
      title: 'label',
      color: AppColors.red,
    );
    ref.watch(labelNotifierProvider.notifier).addLabel(newLabel);
  }

  void _handleRemoveLabel(String labelId) {
    ref.watch(labelNotifierProvider.notifier).removeLabel(labelId);
  }
}

class AddLabelButton extends StatefulWidget {
  final VoidCallback onTap;
  const AddLabelButton({super.key, required this.onTap});

  @override
  State<AddLabelButton> createState() => _AddLabelButtonState();
}

class _AddLabelButtonState extends State<AddLabelButton> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: GestureDetector(
        onTap: widget.onTap,
        child: Container(
          padding: const EdgeInsets.all(5),
          height: 40,
          decoration: BoxDecoration(
            color: AppColors.white15,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: isHovered ? AppColors.white : AppColors.transparent,
            ),
          ),
          child: const Row(
            children: [
              Icon(Icons.add, color: AppColors.white128),

              SizedBox(width: 5),
              Text('Add label', style: AppTextStyles.style3),
            ],
          ),
        ),
      ),
    );
  }
}
